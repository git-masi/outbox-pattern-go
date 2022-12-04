package billing

import (
	"database/sql"
	"git-masi/outbox-pattern-go/cmd/web/events"
	"log"
)

func UpdateCharges(db *sql.DB) events.FulfillmentEventFn {
	return func(event *events.FulfillmentEvent) {
		log.Println("updating charges")
	}
}
