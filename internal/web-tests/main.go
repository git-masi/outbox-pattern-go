package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	spamOrderUpdates()
}

func spamOrderUpdates() {
	client := &http.Client{}
	url := "http://localhost:4000/orders/update"

	for i := 1; i <= 1000; i++ {
		payload, err := json.Marshal(map[string]interface{}{
			"id":     i,
			"status": "shipped",
		})
		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			log.Fatal(err)
		}

		_, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Updated order %d to status %q\n", i, "shipped")

		time.Sleep(2 * time.Second)
	}
}
