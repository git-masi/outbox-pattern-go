package orders

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/alexedwards/flow"
)

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
