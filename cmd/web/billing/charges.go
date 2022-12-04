package billing

import (
	"database/sql"
	"git-masi/outbox-pattern-go/cmd/web/events"
	"log"
	"strconv"
	"strings"
)

func UpdateCharges(db *sql.DB) events.FulfillmentEventFn {
	return func(event *events.FulfillmentEvent) {
		txn, err := db.Begin()
		if err != nil {
			log.Println(err.Error())
			return
		}

		defer txn.Rollback()

		fulfillmentFee, err := getFulfillmentFee(txn, event.Record.EventBody.ClientId)
		if err != nil {
			log.Println("Failed to get fulfillment fee")
			log.Println(err.Error())
			return
		}

		itemFees, err := getItemFees(txn, getItemIds(event.Record.EventBody.ItemIds))
		if err != nil {
			log.Println("Failed to get item fees")
			log.Println(err.Error())
			return
		}

		totalFees := getTotalFees(fulfillmentFee, itemFees)

		err = addNewCharge(txn, totalFees, event.Record.EventBody.OrderId, event.Record.EventBody.ClientId)
		if err != nil {
			log.Println("Failed to add new charge")
			log.Println(err.Error())
			return
		}

		err = txn.Commit()
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
}

func getItemIds(s string) []int {
	strs := strings.Split(s, ",")
	ints := make([]int, len(strs))

	for i, s := range strs {
		ints[i], _ = strconv.Atoi(s)
	}

	return ints
}

func getTotalFees(fulfillmentFee *FulfillmentFee, itemFees []*ItemFee) int {
	total := fulfillmentFee.Rate

	for _, fee := range itemFees {
		total += fee.Rate
	}

	return total
}
