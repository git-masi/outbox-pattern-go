package notifications

import (
	"git-masi/outbox-pattern-go/cmd/web/events"
	"log"
)

// Imagine this was a real notification service which could send email, sms, push notifications, etc.
func NotifyUsersOfFulfillmentEvent(event *events.FulfillmentEvent) {
	log.Printf("notification: %+v\n", *event)
}
