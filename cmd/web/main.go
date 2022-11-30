package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/alexedwards/flow"
)

func main() {
	addr := flag.String("addr", ":4000", "The address for the server to listen on")

	flag.Parse()

	mux := flow.New()

	mux.HandleFunc("/ping", ping, http.MethodGet)

	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
