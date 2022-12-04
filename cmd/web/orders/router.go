package orders

import (
	"database/sql"
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
		mux.HandleFunc(base, getOrders(db), http.MethodGet)

		mux.HandleFunc(addBase("/update/:id"), updateOrder, http.MethodPatch)
	})
}
