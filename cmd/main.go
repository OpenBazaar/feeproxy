package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/OpenBazaar/feeproxy"
	"github.com/gocraft/health"
)

var cache string

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, cache)
}

func run() {
	// Create instrumentation stream
	stream := health.NewStream()
	stream.AddSink(&health.WriterSink{os.Stdout})
	stream.Event("starting")

	// Create an update ticker and update fees on each tick
	t := time.NewTicker(time.Second * 30)
	for ; true; <-t.C {

		// Query the upstream API
		job := stream.NewJob("query")
		newDataJSON, err := feeproxy.Query()
		if err == nil {
			cache = string(newDataJSON)
			job.Complete(health.Success)
			continue
		}

		job.EventErr("query", err)
		job.Complete(health.Error)
	}
}

func main() {
	// Start querying the upstream API
	go run()

	// Get listening interface address
	addr := getListenAddr()
	fmt.Println("Starting server on", addr)

	// Listen for requests
	http.HandleFunc("/", handler)
	http.ListenAndServe(addr, nil)
}

func getListenAddr() string {
	val := os.Getenv("FEEPROXY_INTERFACE")
	if val == "" {
		return ":8080"
	}

	return val
}
