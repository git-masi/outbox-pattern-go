package orders

import (
	"database/sql"
	"encoding/json"
	"errors"
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
			return
		}

		w.Write(res)
	}
}

func updateOrder(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			http.Error(w, "Content type is not application/json", http.StatusBadRequest)
			return
		}

		var update struct {
			OrderId int    `json:"id"`
			Status  string `json:"status"`
		}
		var unmarshalErr *json.UnmarshalTypeError

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&update)
		if err != nil {
			if errors.As(err, &unmarshalErr) {
				http.Error(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
			} else {
				http.Error(w, "Bad Request "+err.Error(), http.StatusBadRequest)
			}
			return
		}

		txn, err := db.Begin()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer txn.Rollback()

		err = updateOrderStatus(txn, update.OrderId, update.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add an event to the "outbox"
		if update.Status == "shipped" {
			fd, err := getFulfillmentData(txn, update.OrderId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			data, err := json.Marshal(fd)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = addFulfillmentStatusEvent(txn, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		err = txn.Commit()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("OK"))
	}
}
