package providerimpl

import (
	`stripe.com/docs/payments/core/api`
	`stripe.com/docs/payments/core/model`
	`stripe.com/docs/payments/provider`
)

type razorpayProvider struct {
}

func (r *razorpayProvider) Create(payment *model.Payment) (*api.PaymentResponse, error) {
	return nil, nil
}

func NewRazorpayProvider() provider.Payment {
	return &razorpayProvider{}
}
