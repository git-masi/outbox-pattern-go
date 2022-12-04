package orders

import (
	"database/sql"
	"git-masi/outbox-pattern-go/cmd/web/events"
	"time"
)

type Order struct {
	Id          int    `json:"id"`
	Created     string `json:"created"`
	LastUpdated string `json:"lastUpdated"`
	Status      string `json:"status"`
	ClientId    string `json:"clientId"`
}

func readAllOrders(db *sql.DB) ([]*Order, error) {
	// In a real application we should have some limits on the number of orders
	// returned from the Db.
	// It would also be good to add optional conditions like orders that were
	// created after a certain date or with a specific status or client ID.
	stmt := `SELECT * FROM orders;`

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

func updateOrderStatus(txn *sql.Tx, orderId int, status string) error {
	stmt := `UPDATE orders SET "status" = $1 WHERE orders.id = $2`

	_, err := txn.Exec(stmt, status, orderId)
	if err != nil {
		return err
	}

	return nil
}

func addFulfillmentStatusEvent(txn *sql.Tx, body []byte) error {
	now := time.Now()
	iso := now.Format(time.RFC3339)
	stmt := `INSERT INTO order_fulfillment_events(created, event_body) VALUES($1, $2)`

	_, err := txn.Exec(stmt, iso, body)
	if err != nil {
		return err
	}

	return nil
}

func getFulfillmentData(txn *sql.Tx, orderId int) (*events.FulfillmentEventBody, error) {
	stmt := `
		SELECT
			o.id AS order_id,
			o.client_id,
			array_to_string(array_agg(i.id), ',') AS item_ids
		FROM orders o
		JOIN order_items oi ON o.id = oi.order_id
		JOIN items i ON oi.item_id = i.id
		WHERE o.id = $1
		GROUP BY o.id;
	`

	fd := &events.FulfillmentEventBody{}

	row := txn.QueryRow(stmt, orderId)

	err := row.Scan(&fd.OrderId, &fd.ClientId, &fd.ItemIds)
	if err != nil {
		return nil, err
	}

	return fd, nil
}
