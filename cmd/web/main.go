package main

import (
	"database/sql"
	"encoding/json"
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

	OrderRouter(mux, db)

	err = http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}

type Order struct {
	Id          int    `json:"id"`
	Created     string `json:"created"`
	LastUpdated string `json:"lastUpdated"`
	Status      string `json:"status"`
	ClientId    string `json:"clientId"`
}

func OrderRouter(mux *flow.Mux, db *sql.DB) {
	base := "/orders"
	addBase := func(path string) string {
		return base + path
	}

	// Using a group in case any custom middleware is needed in the future
	mux.Group(func(m *flow.Mux) {
		mux.HandleFunc(base, readAllOrders(db), http.MethodGet)

		mux.HandleFunc(addBase("/update"), updateOrder, http.MethodPatch)
	})
}

func readAllOrders(db *sql.DB) http.HandlerFunc {
	readOrders := func() ([]*Order, error) {
		stmt := `select * from orders;`

		rows, err := db.Query(stmt)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		orders := []*Order{}

		for rows.Next() {
			o := &Order{}

			err = rows.Scan(&o.Id, &o.Created, &o.LastUpdated, &o.Status, &o.ClientId)
			if err != nil {
				return nil, err
			}

			orders = append(orders, o)
		}

		return orders, nil
	}

	return func(w http.ResponseWriter, r *http.Request) {
		orders, err := readOrders()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")

		res, err := json.Marshal(orders)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write(res)
	}
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update!"))
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
