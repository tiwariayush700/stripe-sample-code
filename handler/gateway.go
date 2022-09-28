package handler

import (
	`context`
	`encoding/json`
	`io/ioutil`
	`log`
	`net/http`
	`strings`
	`time`

	`github.com/gin-gonic/gin`
	"github.com/stripe/stripe-go/v71/webhook"
	`github.com/stripe/stripe-go/v72`

	`stripe.com/docs/payments/core/api`
	`stripe.com/docs/payments/core/config`
	`stripe.com/docs/payments/core/model`
	`stripe.com/docs/payments/core/util`
	`stripe.com/docs/payments/service`
)

type Gateway struct {
	PGService     service.Gateway
	PaymentConfig config.Payment
	//payoutService
	//refundService etc......
}

func NewGateway(gatewayService service.Gateway, paymentConfig config.Payment) *Gateway {
	return &Gateway{PGService: gatewayService, PaymentConfig: paymentConfig}
}

func (p *Gateway) RegisterRoutes(ctx context.Context, router *gin.Engine) {

	paymentRouterGroup := router.Group("/payment")
	{
		paymentRouterGroup.POST("", p.Create(ctx))
	}
	webhookRouterGroup := router.Group("/webhook")
	{
		webhookRouterGroup.POST("", p.ListenWebhook(ctx))
	}
}

func (p *Gateway) Create(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {

		var paymentRequest api.PaymentRequest
		err := c.ShouldBind(&paymentRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, api.NewHTTPError(api.ErrorCodeInvalidRequestPayload, "Invalid request body"))
			return
		}

		paymentProvider := p.PGService.GetPaymentProvider(paymentRequest.Currency)

		total := model.NewMoney(paymentRequest.Amount, paymentRequest.Currency)
		payment := model.NewPayment(&total)
		paymentResponse, err := paymentProvider.Create(payment)
		if err != nil {
			c.JSON(http.StatusBadRequest, api.NewHTTPError(api.ErrorCodePaymentGateway, "Failed to initiate payment"))
			return
		}

		payment.ExternalID = paymentResponse.ID
		err = p.PGService.CreatePayment(ctx, payment)
		if err != nil {
			c.JSON(http.StatusBadRequest, api.NewHTTPError(api.ErrorCodeDatabaseFailure, "Failed to save payment"))
			return
		}

		c.JSON(http.StatusOK, paymentResponse)

	}
}

func (p *Gateway) ListenWebhook(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {

		//decoding request
		//invalid request should be non 200
		const MaxBodyBytes = int64(65536)
		payload, err := ioutil.ReadAll(http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes))
		if err != nil {
			log.Printf("ReceiveStripeConfirmPayment binding error : %v", err)
			c.JSON(http.StatusBadRequest, err)
			return
		}

		stripe.Key = p.PaymentConfig.Secret
		_, err = webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"), p.PaymentConfig.WebhookSecret)
		if err != nil {
			log.Printf("ReceiveStripeConfirmPayment webhook signature validation error : %v", err)
			c.JSON(http.StatusBadRequest, err)
			return
		}

		event := &stripe.Event{}
		err = json.Unmarshal(payload, event)
		if err != nil {
			log.Printf("ReceiveStripeConfirmPayment unmarshalling error : %v", err)
			c.JSON(http.StatusBadRequest, err)
			return
		}

		log.Printf("ReceiveStripeConfirmPayment : Webhook event from stripe : %+v", event)

		//switch to different event listeners based on event.Type
		//event.Type = "transfer.created"
		eventTypes := strings.Split(event.Type, ".")
		var eventType string
		if len(eventTypes) > 0 {
			eventType = eventTypes[0]
		}
		switch eventType {
		case "payment_intent":
			//switch to payment listener

			//TODO add impl
			_ = p.PGService.PaymentEventListener(ctx, event)

			paymentIntent := &stripe.PaymentIntent{}
			err = json.Unmarshal(event.Data.Raw, &paymentIntent)
			if err != nil {
				log.Printf("ReceiveStripeConfirmPayment : err unmarshalling paymentIntent with event : %v err : %v", event, err)
				c.JSON(http.StatusUnprocessableEntity, err)
				return
			}

			var payment *model.Payment
			if paymentID, ok := paymentIntent.Metadata["paymentId"]; ok {
				payment, err = p.PGService.GetPayment(ctx, paymentID)
				if err != nil {
					//if your db is down, etc you might want a retry so non 200 status
					c.JSON(http.StatusBadRequest, err)
					return
				}

				log.Printf("Payment %+v", *payment)

				if util.IsDuplicateEvent(event.Request.IdempotencyKey, payment.WebhookEvents) {
					//since a duplicate event
					//return 200, so it's not triggered again
					log.Printf("ReceiveStripeConfirmPayment : duplicate event")
					c.JSON(http.StatusOK, "Duplicate event")
					return
				}

				if len(payment.WebhookEvents) == 0 {
					payment.WebhookEvents = []model.WebhookEvent{{
						ID:          event.Request.IdempotencyKey,
						Type:        model.WebhookEventType(event.Type),
						ProcessedAt: time.Now(),
					}}
				} else {
					payment.WebhookEvents = append(payment.WebhookEvents, model.WebhookEvent{
						ID:          event.Request.IdempotencyKey,
						Type:        model.WebhookEventType(event.Type),
						ProcessedAt: time.Now(),
					})
				}

				if v, ok := model.StripeStatusMapper[paymentIntent.Status]; ok {
					payment.Status = v
				}

				err = p.PGService.UpdatePayment(ctx, payment)
				if err != nil {
					c.JSON(http.StatusBadRequest, err)
					return
				}
				c.JSON(http.StatusOK, payment)
			}
		case "payout":
			//switch to payout listener

		}

	}
}
