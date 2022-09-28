package util

import (
	`stripe.com/docs/payments/core/model`
)

func IsDuplicateEvent(eventID string, events []model.WebhookEvent) bool {
	for _, event := range events {
		if eventID == event.ID {
			return true
		}
	}

	return false
}
