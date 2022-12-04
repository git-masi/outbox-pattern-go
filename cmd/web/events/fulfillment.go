package events

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
