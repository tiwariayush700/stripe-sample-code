package model

import (
	`time`

	`github.com/hashicorp/go-uuid`
	`github.com/stripe/stripe-go/v72`
)

type Payment struct {
	ID            string         `json:"id" bson:"_id"`
	Total         *Money         `json:"amount" bson:"amount"`
	Status        Status         `json:"status" bson:"status"`
	CreatedAt     time.Time      `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" bson:"updated_at"`
	WebhookEvents []WebhookEvent `json:"webhook_events" bson:"webhookEvents"`
	ExternalID    string         `json:"external_id" bson:"external_id"`
}

type Status string

const (
	SuccessStatus = Status("success")
	FailedStatus  = Status("failed")
	PendingStatus = Status("pending")
)

var StripeStatusMapper = map[stripe.PaymentIntentStatus]Status{
	stripe.PaymentIntentStatusSucceeded: SuccessStatus,
	stripe.PaymentIntentStatusCanceled:  FailedStatus,
}

func NewPayment(total *Money) *Payment {
	id, _ := uuid.GenerateUUID()
	return &Payment{
		ID:            id,
		Total:         total,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		WebhookEvents: nil,
	}
}

type WebhookEvent struct {
	ID          string           `json:"id" bson:"_id"`
	Type        WebhookEventType `json:"type" bson:"type"`
	ProcessedAt time.Time        `json:"processed_at" bson:"processed_at"`
}

type WebhookEventType string

const (
	//stripe
	WebhookEventTypePaymentIntentSucceeded = WebhookEventType("payment_intent.succeeded")
	WebhookEventTypePaymentIntentFailed    = WebhookEventType("payment_intent.payment_failed")
	WebhookEventTypePaymentIntentCanceled  = WebhookEventType("payment_intent.canceled")

	//razorpay
	WebhookEventTypePaymentCaptured = WebhookEventType("payment.captured")
	WebhookEventTypePaymentFailed   = WebhookEventType("payment.failed")
)
