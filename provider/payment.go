package provider

import (
	`stripe.com/docs/payments/core/api`
	`stripe.com/docs/payments/core/model`
)

type Payment interface {
	Create(payment *model.Payment) (*api.PaymentResponse, error)
}
