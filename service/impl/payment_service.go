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

type gatewayServiceImpl struct {
	PaymentConfig     config.Payment
	PaymentRepository repository.Payment
}

func (p *gatewayServiceImpl) UpdatePayment(ctx context.Context, payment *model.Payment) error {
	return p.PaymentRepository.Update(ctx, payment)
}

func (p *gatewayServiceImpl) PaymentEventListener(ctx context.Context, data interface{}) error {
	//Based on the currency switch to the particular providers impl
	//Move all the status handling code here
	return nil
}

func (p *gatewayServiceImpl) GetPayment(ctx context.Context, id string) (*model.Payment, error) {
	return p.PaymentRepository.Get(ctx, id)
}

func (p *gatewayServiceImpl) CreatePayment(ctx context.Context, payment *model.Payment) error {
	payment.Status = model.PendingStatus
	return p.PaymentRepository.Create(ctx, payment)
}

func (p *gatewayServiceImpl) GetPaymentProvider(currency string) provider.Payment {
	switch currency {
	case "USD":
		return providerimpl.NewStripeProvider(p.PaymentConfig)
	}

	return nil
}

func NewGatewayServiceImpl(paymentConfig config.Payment, paymentRepository repository.Payment) service.Gateway {
	return &gatewayServiceImpl{PaymentConfig: paymentConfig, PaymentRepository: paymentRepository}
}
