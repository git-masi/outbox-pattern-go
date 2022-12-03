package main

import (
	"database/sql"
	"flag"
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

	go listenForFulfillmentEvent(*dsn, []func(*pq.Notification){NotifyUsersOfFulfillmentEvent, UpdateCharges(db)})

	mux := flow.New()

	mux.HandleFunc("/ping", ping, http.MethodGet)

	err = http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func listenForFulfillmentEvent(dsn string, subscribers []func(*pq.Notification)) {
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
		case msg := <-listener.Notify:
			for _, fn := range subscribers {
				go fn(msg)
			}
		case <-time.After(90 * time.Second):
			go listener.Ping()
			log.Println("No new notifications in past 90 seconds, pinging DB to ensure connection is still alive")
		}
	}
}

// Imagine this was a real notification service which could send email, sms, push notifications, etc.
func NotifyUsersOfFulfillmentEvent(notification *pq.Notification) {
	log.Printf("notification: %+v\n", notification)
}

func UpdateCharges(db *sql.DB) func(notification *pq.Notification) {
	return func(notification *pq.Notification) {
		log.Println("updating charges")
	}
}
