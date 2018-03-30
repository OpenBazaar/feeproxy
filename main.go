package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gocraft/health"
)

const api string = "https://bitcoinfees.21.co/api/v1/fees/list"

type Fees struct {
	Fees []FeeLevel `json:"fees"`
}

type FeeLevel struct {
	MinFee     int `json:"minFee"`
	MaxFee     int `json:"maxFee"`
	DayCount   int `json:"dayCount"`
	MemCount   int `json:"memCount"`
	MinDelay   int `json:"minDelay"`
	MaxDelay   int `json:"maxDelay"`
	MaxMinutes int `json:"maxMinutes"`
}

type FeeCache struct {
	Priority int `json:"priority"`
	Normal   int `json:"normal"`
	Economic int `json:"economic"`
}

var cache string
var httpClient http.Client = http.Client{Timeout: time.Second * 10}

func Query() error {
	resp, err := httpClient.Get(api)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(resp.Body)
	fees := new(Fees)
	err = decoder.Decode(fees)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	low := FeeLevel{MaxDelay: 0}
	medium := FeeLevel{MaxDelay: 0}
	high := FeeLevel{MaxDelay: 10}
	for _, fl := range fees.Fees {
		if (fl.MaxDelay > low.MaxDelay && fl.MaxDelay <= 6) || (fl.MaxDelay == low.MaxDelay && fl.MaxFee < low.MaxFee) {
			low = fl
		}
		if (fl.MaxDelay > medium.MaxDelay && fl.MaxDelay <= 3) || (fl.MaxDelay == medium.MaxDelay && fl.MaxFee < medium.MaxFee) {
			medium = fl
		}
		if (fl.MaxDelay < high.MaxDelay) || (fl.MaxDelay == high.MaxDelay && fl.MaxFee < high.MaxFee) {
			high = fl
		}
	}
	feeCache := FeeCache{
		Priority: (high.MaxFee + high.MinFee) / 2,
		Normal:   (medium.MaxFee + medium.MinFee) / 2,
		Economic: (low.MaxFee + low.MinFee) / 2,
	}

	if feeCache.Priority == 0 {
		feeCache.Priority = 1
	}
	if feeCache.Normal == 0 {
		feeCache.Normal = 1
	}
	if feeCache.Economic == 0 {
		feeCache.Economic = 1
	}

	out, err := json.Marshal(&feeCache)
	if err != nil {
		return err
	}
	cache = string(out)
	return nil
}

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
		err := Query()
		if err == nil {
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
