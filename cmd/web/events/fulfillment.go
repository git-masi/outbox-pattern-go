package events

import (
	"encoding/json"
	"log"
	"time"

	"github.com/lib/pq"
)

type FulfillmentEventFn func(*FulfillmentEvent)

type FulfillmentEvent struct {
	Operation string            `json:"operation"`
	Record    FulfillmentRecord `json:"record"`
}

type FulfillmentRecord struct {
	Id        int                  `json:"id"`
	Created   string               `json:"created"`
	EventBody FulfillmentEventBody `json:"event_body"`
}

type FulfillmentEventBody struct {
	OrderId  int    `json:"order_id"`
	ClientId int    `json:"client_id"`
	ItemIds  string `json:"item_ids"`
}

// This is effectively a pubsub
func ListenForFulfillmentEvent(dsn string, subscribers []FulfillmentEventFn) {
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Println(err.Error())
		}
	}

	minReconn := 10 * time.Second
	maxReconn := time.Minute
	listener := pq.NewListener(dsn, minReconn, maxReconn, reportProblem)

	err := listener.Listen("fulfillment_event")
	if err != nil {
		log.Println(err.Error())
		return
	}

	for {
		select {
		case ntf := <-listener.Notify:
			var event FulfillmentEvent

			err := json.Unmarshal([]byte(ntf.Extra), &event)
			if err != nil {
				log.Println(err.Error())
				return
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
