package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"git-masi/outbox-pattern-go/cmd/web/events"
	"git-masi/outbox-pattern-go/cmd/web/notifications"
	"git-masi/outbox-pattern-go/cmd/web/orders"
	"git-masi/outbox-pattern-go/internal/db"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/flow"
	"github.com/lib/pq"
)

func main() {
	addr := flag.String("addr", ":4000", "The address for the server to listen on")
	dsn := flag.String("dsn", "postgres://postgres:postgres@localhost/outbox?sslmode=disable", "PostgreSQL DSN")

	flag.Parse()

	db, err := db.OpenDb(*dsn)
	if err != nil {
		log.Fatal(err)
	}

	go listenForFulfillmentEvent(*dsn, []events.FulfillmentEventFn{notifications.NotifyUsersOfFulfillmentEvent, UpdateCharges(db)})

	mux := flow.New()

	orders.OrderRouter(mux, db)

	err = http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}

func listenForFulfillmentEvent(dsn string, subscribers []events.FulfillmentEventFn) {
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Println(err.Error())
		}
	}

	minReconn := 10 * time.Second
	maxReconn := time.Minute
	listener := pq.NewListener(dsn, minReconn, maxReconn, reportProblem)
	// Listen for a notification from the `fulfillment_event` channel
	err := listener.Listen("fulfillment_event")
	if err != nil {
		panic(err)
	}

	for {
		select {
		case ntf := <-listener.Notify:
			var event events.FulfillmentEvent

			err := json.Unmarshal([]byte(ntf.Extra), &event)
			if err != nil {
				log.Println(err.Error())
			}

			for _, fn := range subscribers {
				go fn(&event)
			}
		case <-time.After(90 * time.Second):
			go listener.Ping()
			log.Println("No new notifications in past 90 seconds, pinging DB to ensure connection is still alive")
		}
	}
}

func UpdateCharges(db *sql.DB) events.FulfillmentEventFn {
	return func(event *events.FulfillmentEvent) {
		log.Println("updating charges")
	}
}
