package service

import (
	`context`

	`stripe.com/docs/payments/core/model`
	`stripe.com/docs/payments/provider`
)

type Payment interface {
	GetPaymentProvider(currency string) provider.Payment
	CreatePayment(ctx context.Context, payment *model.Payment) error
}
