package main

import (
	"context"
	"database/sql"
	"flag"
	"git-masi/outbox-pattern-go/internal/functional"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/icrowley/fake"
	"github.com/lib/pq"
)

type Client struct {
	Created     string
	LastUpdated string
	Name        string
}

func main() {
	dsn := flag.String("dsn", "postgres://postgres:postgres@localhost/outbox?sslmode=disable", "PostgreSQL DSN")

	flag.Parse()

	db, err := openDb(*dsn)
	if err != nil {
		log.Fatal(err)
	}

	createFakeData(db)
}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createFakeData(db *sql.DB) {
	now := time.Now()
	iso := now.Format(time.RFC3339)
	orderStatus := "created"
	numCompanies := 5
	numOrders := 1000
	numItems := 100
	maxOrderItems := 5

	rand.Seed(now.UnixMicro())

	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	test := true
	if !test {
		createClients(txn, numCompanies, iso)
		createOrders(txn, numOrders, numCompanies, orderStatus, iso)
	}

	clientItems := createItems(txn, numItems, numCompanies, iso)
	createOrderItems(txn, clientItems, numCompanies, numOrders, maxOrderItems)

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func createClients(txn *sql.Tx, numCompanies int, iso string) {
	stmt, err := txn.Prepare(pq.CopyIn("clients", "created", "last_updated", "name"))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < numCompanies; i++ {
		_, err = stmt.Exec(iso, iso, fake.Company())
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func createOrders(txn *sql.Tx, numOrders int, numCompanies int, status string, iso string) {
	stmt, err := txn.Prepare(pq.CopyIn("orders", "created", "last_updated", "status", "client_id"))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < numOrders; i++ {
		_, err = stmt.Exec(iso, iso, status, rand.Intn(numCompanies)+1)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func createItems(txn *sql.Tx, numItems int, numCompanies int, iso string) map[int][]int {
	stmt, err := txn.Prepare(pq.CopyIn("items", "created", "last_updated", "name", "description", "price", "client_id"))
	if err != nil {
		log.Fatal(err)
	}

	clientItems := map[int][]int{}

	// Create items, `i` will be the item ID
	for i := 1; i <= numItems; i++ {
		price, err := strconv.Atoi(fake.DigitsN(4))
		if err != nil {
			log.Fatal(err)
		}

		companyId := rand.Intn(numCompanies) + 1

		_, err = stmt.Exec(iso, iso, fake.ProductName(), fake.Sentence(), price, companyId)
		if err != nil {
			log.Fatal(err)
		}

		if _, ok := clientItems[companyId]; !ok {
			clientItems[companyId] = []int{}
		}

		clientItems[companyId] = append(clientItems[companyId], i)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	return clientItems
}

func createOrderItems(txn *sql.Tx, clientItems map[int][]int, numCompanies, numOrders, maxOrderItems int) {
	stmt, err := txn.Prepare(pq.CopyIn("order_items", "order_id", "item_id"))
	if err != nil {
		log.Fatal(err)
	}

	// Create order items for each order
	for orderId := 1; orderId <= numOrders; orderId++ {
		companyId := rand.Intn(numCompanies) + 1
		availableItems := clientItems[companyId]
		// Number of items for this order
		numOrderItems := rand.Intn(maxOrderItems) + 1

		if numOrderItems > len(availableItems) {
			numOrderItems = len(availableItems)
		}

		for j := 0; j < numOrderItems; j++ {
			itemIdx := rand.Intn(len(availableItems))
			itemId := availableItems[itemIdx]
			availableItems = functional.Filter(availableItems, func(id int) bool { return id != itemId })

			_, err = stmt.Exec(orderId, itemId)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func createFulfillmentFees(txn *sql.Tx, numCompanies int, iso string) {
	stmt, err := txn.Prepare(pq.CopyIn("fulfillment_fees", "created", "last_updated", "rate", "client_id"))
	if err != nil {
		log.Fatal(err)
	}

	for clientId := 1; clientId <= numCompanies; clientId++ {
		rate := rand.Intn(10) + 1

		_, err = stmt.Exec(iso, iso, rate, clientId)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func createItemFees(txn *sql.Tx, clientItems map[int][]int, iso string) {
	stmt, err := txn.Prepare(pq.CopyIn("item_fees", "created", "last_updated", "rate", "fee_description", "item_id", "client_id"))
	if err != nil {
		log.Fatal(err)
	}

	itemFeeDescriptions := []string{"heavy item", "large item", "fragile item"}

	for clientId, items := range clientItems {
		// Most items should not have fees
		// Only 1 in 6 clients should have items with fees
		if rand.Intn(6) > 4 {
			// Randomly select a portion of the items to add fees for
			numItemsWithFees := rand.Intn(len(items)) + 1
			availableItems := items

			for i := 0; i < numItemsWithFees; i++ {
				itemIdx := rand.Intn(len(availableItems))
				itemId := availableItems[itemIdx]
				availableItems = functional.Filter(availableItems, func(id int) bool { return id != itemId })
				rate := rand.Intn(10) + 1
				feeDescription := itemFeeDescriptions[rand.Intn(len(itemFeeDescriptions))]

				_, err = stmt.Exec(iso, iso, rate, feeDescription, itemId, clientId)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
}
