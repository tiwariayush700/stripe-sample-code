package repository

import (
	`context`

	`stripe.com/docs/payments/core/model`
)

type Payment interface {
	Create(ctx context.Context, payment *model.Payment) error
	Get(ctx context.Context, id string) (*model.Payment, error)
	Update(ctx context.Context, payment *model.Payment) error
}
