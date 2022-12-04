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
