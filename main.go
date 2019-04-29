package feeproxy

import (
	"encoding/json"
	"net/http"
	"time"
)

const api = "https://bitcoinfees.earn.com/api/v1/fees/recommended"

type sourceData struct {
	Priority int `json:"fastestFee"`
	Normal   int `json:"halfHourFee"`
	Economic int `json:"hourFee"`
}

type responseData struct {
	Priority int `json:"priority"`
	Normal   int `json:"normal"`
	Economic int `json:"economic"`
}

var httpClient http.Client = http.Client{Timeout: time.Second * 30}

func Query() ([]byte, error) {
	resp, err := httpClient.Get(api)
	if err != nil {
		return nil, err
	}

	feeData := &sourceData{}
	err = json.NewDecoder(resp.Body).Decode(feeData)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	return json.Marshal(&responseData{
		Priority: feeData.Priority,
		Normal:   feeData.Normal,
		Economic: feeData.Economic,
	})
}
