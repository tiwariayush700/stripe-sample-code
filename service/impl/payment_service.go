package serviceimpl

import (
	`context`

	`stripe.com/docs/payments/core/config`
	`stripe.com/docs/payments/core/model`
	`stripe.com/docs/payments/provider`
	providerimpl `stripe.com/docs/payments/provider/impl`
	`stripe.com/docs/payments/repository`
	`stripe.com/docs/payments/service`
)

type paymentServiceImpl struct {
	PaymentConfig     config.Payment
	PaymentRepository repository.Payment
}

func (p *paymentServiceImpl) CreatePayment(ctx context.Context, payment *model.Payment) error {
	return p.PaymentRepository.Create(ctx, payment)
}

func (p *paymentServiceImpl) GetPaymentProvider(currency string) provider.Payment {
	switch currency {
	case "USD":
		return providerimpl.NewStripeProvider(p.PaymentConfig)
	}

	return nil
}

func NewPaymentServiceImpl(paymentConfig config.Payment, paymentRepository repository.Payment) service.Payment {
	return &paymentServiceImpl{PaymentConfig: paymentConfig, PaymentRepository: paymentRepository}
}
