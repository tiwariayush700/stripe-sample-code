package handler

import (
	`context`
	`net/http`

	`github.com/gin-gonic/gin`

	`stripe.com/docs/payments/core/api`
	`stripe.com/docs/payments/core/model`
	`stripe.com/docs/payments/service`
)

type Payment struct {
	PaymentService service.Payment
}

func NewPayment(paymentService service.Payment) *Payment {
	return &Payment{PaymentService: paymentService}
}

func (p *Payment) RegisterRoutes(ctx context.Context, router *gin.Engine) {

	paymentRouterGroup := router.Group("/payment")
	{
		paymentRouterGroup.POST("", p.Create(ctx))
	}
}

func (p *Payment) Create(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {

		var paymentRequest api.PaymentRequest
		err := c.ShouldBind(&paymentRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, api.NewHTTPError(api.ErrorCodeInvalidRequestPayload, "Invalid request body"))
			return
		}

		paymentProvider := p.PaymentService.GetPaymentProvider(paymentRequest.Currency)

		total := model.NewMoney(paymentRequest.Amount, paymentRequest.Currency)
		payment := model.NewPayment(&total)
		paymentResponse, err := paymentProvider.Create(payment)
		if err != nil {
			c.JSON(http.StatusBadRequest, api.NewHTTPError(api.ErrorCodePaymentGateway, "Failed to initiate payment"))
			return
		}

		err = p.PaymentService.CreatePayment(ctx, payment)
		if err != nil {
			c.JSON(http.StatusBadRequest, api.NewHTTPError(api.ErrorCodeDatabaseFailure, "Failed to save payment"))
			return
		}

		c.JSON(http.StatusOK, paymentResponse)

	}
}
