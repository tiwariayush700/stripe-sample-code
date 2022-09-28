package service

import (
	`context`

	`stripe.com/docs/payments/core/model`
	`stripe.com/docs/payments/provider`
)

type Gateway interface {
	GetPaymentProvider(currency string) provider.Payment
	CreatePayment(ctx context.Context, payment *model.Payment) error
	GetPayment(ctx context.Context, id string) (*model.Payment, error)
	UpdatePayment(ctx context.Context, payment *model.Payment) error
	PaymentEventListener(ctx context.Context, data interface{}) error
}
