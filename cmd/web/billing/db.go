package billing

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type FulfillmentFee struct {
	Rate     int `json:"rate"`
	ClientId int `json:"client_id"`
}

type ItemFee struct {
	Id             int    `json:"-"`
	Created        string `json:"-"`
	LastUpdated    string `json:"-"`
	Rate           int    `json:"rate"`
	FeeDescription string `json:"fee_description"`
	ItemId         int    `json:"item_id"`
	ClientId       int    `json:"client_id"`
}

func getFulfillmentFee(txn *sql.Tx, clientId int) (*FulfillmentFee, error) {
	stmt := `SELECT ff.rate FROM fulfillment_fees ff WHERE ff.client_id = $1`
	fee := &FulfillmentFee{}
	row := txn.QueryRow(stmt, clientId)

	err := row.Scan(&fee.Rate)
	if err != nil {
		return nil, err
	}

	return fee, nil
}

func getItemFees(txn *sql.Tx, itemIds []int) ([]*ItemFee, error) {
	stmt := `SELECT i.* FROM item_fees i WHERE i.item_id = ANY($1)`
	fees := []*ItemFee{}

	rows, err := txn.Query(stmt, pq.Array(itemIds))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		fee := &ItemFee{}

		err = rows.Scan(&fee.Rate, &fee.FeeDescription, &fee.ItemId, &fee.ClientId)
		if err != nil {
			return nil, err
		}

		fees = append(fees, fee)
	}

	return fees, nil
}

func addNewCharge(txn *sql.Tx, amount, orderId, clientId int) error {
	now := time.Now()
	iso := now.Format(time.RFC3339)
	stmt := `INSERT INTO charges(created, amount, order_id, client_id) VALUES($1,$2,$3,$4)`

	_, err := txn.Exec(stmt, iso, amount, orderId, clientId)
	if err != nil {
		return err
	}

	return nil
}
