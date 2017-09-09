package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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
var lastSuccessfulQuery time.Time

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
	out, err := json.Marshal(&feeCache)
	if err != nil {
		return err
	}
	cache = string(out)
	lastSuccessfulQuery = time.Now()
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, cache)
}

func run() {
	t := time.NewTicker(time.Second * 30)
	for ; true; <-t.C {
		err := Query()
		if err != nil && time.Since(lastSuccessfulQuery) > time.Minute*10 {
			// Send a notification or something
		}
	}
}

func main() {
	go run()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
