package notifications

import (
	"log"

	"github.com/lib/pq"
)

// Imagine this was a real notification service which could send email, sms, push notifications, etc.
func NotifyUsersOfFulfillmentEvent(notification *pq.Notification) {
	log.Printf("notification: %+v\n", notification)
}
