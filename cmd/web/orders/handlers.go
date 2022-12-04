package orders

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func getOrders(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orders, err := readAllOrders(db)

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
