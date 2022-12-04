package orders

import "database/sql"

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

func updateOrderStatus(db *sql.DB, orderId int, status string) error {
	stmt := `UPDATE orders SET "status" = $1 WHERE orders.id = $2`

	_, err := db.Exec(stmt, status, orderId)
	if err != nil {
		return err
	}

	return nil
}
