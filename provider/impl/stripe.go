package providerimpl

import (
	`log`
	`strings`

	`github.com/stripe/stripe-go/v72`
	`github.com/stripe/stripe-go/v72/paymentintent`

	`stripe.com/docs/payments/core/api`
	`stripe.com/docs/payments/core/config`
	`stripe.com/docs/payments/core/model`
	`stripe.com/docs/payments/provider`
)

type stripeProvider struct {
	PaymentConfig config.Payment
}

func (s *stripeProvider) Create(payment *model.Payment) (*api.PaymentResponse, error) {
	paymentConfig := s.PaymentConfig

	log.Println("sec == ", paymentConfig.Secret)
	stripe.Key = paymentConfig.Secret //"sk_test_****"

	stripeParams := stripe.Params{Metadata: map[string]string{
		"paymentId": payment.ID,
		//"anyRandomMetaData": "CB Payments",
	}}

	sParams := &stripe.PaymentIntentParams{
		Params: stripeParams,
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Amount:   stripe.Int64(payment.Total.Amount()),
		Currency: stripe.String(strings.ToLower(payment.Total.Currency().Code)),
		//AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
		//	Enabled: stripe.Bool(true),
		//},
		//for splitting money instantly while creating payment
		//ApplicationFeeAmount: stripe.Int64(123),
		//TransferData: &stripe.PaymentIntentTransferDataParams{
		//	Destination: stripe.String(p.accountId),
		//},
	}

	sParams.SetStripeAccount(paymentConfig.AccountID)

	log.Printf("Stripe PaymentIntentParams => %v", stripeParams)

	pi, err := paymentintent.New(sParams)
	if err != nil {

		log.Printf("paymentintent.New : err %v", err)

		//stripeErr, _ := err.(*stripe.Error)
		//status := stripeErr.HTTPStatusCode
		return nil, err
	}

	response := &api.PaymentResponse{
		ID:           pi.ID,
		PublicKey:    paymentConfig.Key,
		ClientSecret: pi.ClientSecret,
		Extras:       nil,
	}

	return response, nil
}

func NewStripeProvider(paymentConfig config.Payment) provider.Payment {
	return &stripeProvider{PaymentConfig: paymentConfig}
}
