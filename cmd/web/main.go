package main

import (
	"flag"
	"git-masi/outbox-pattern-go/cmd/web/billing"
	"git-masi/outbox-pattern-go/cmd/web/events"
	"git-masi/outbox-pattern-go/cmd/web/notifications"
	"git-masi/outbox-pattern-go/cmd/web/orders"
	"git-masi/outbox-pattern-go/internal/db"
	"log"
	"net/http"

	"github.com/alexedwards/flow"
	_ "github.com/lib/pq"
)

func main() {
	addr := flag.String("addr", ":4000", "The address for the server to listen on")
	dsn := flag.String("dsn", "postgres://postgres:postgres@localhost/outbox?sslmode=disable", "PostgreSQL DSN")

	flag.Parse()

	db, err := db.OpenDb(*dsn)
	if err != nil {
		log.Fatal(err)
	}

	// A real pubsub should be able to subscribe and unsubscribe but for simplicity
	// we can init everything here to focus on functionality
	go events.ListenForFulfillmentEvent(*dsn, []events.FulfillmentEventFn{notifications.NotifyUsersOfFulfillmentEvent, billing.UpdateCharges(db)})

	mux := flow.New()

	orders.OrderRouter(mux, db)

	err = http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
