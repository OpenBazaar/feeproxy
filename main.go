package feeproxy

import (
	"encoding/json"
	"net/http"
	"time"
)

const api = "https://bitcoinfees.earn.com/api/v1/fees/recommended"
const bcinfoAPI = "https://api.blockchain.info/mempool/fees"

type sourceData struct {
	Priority int `json:"fastestFee"`
	Normal   int `json:"halfHourFee"`
	Economic int `json:"hourFee"`
}

type limits struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type sourceDataBCI struct {
	Limits   limits `json:"limits"`
	Regular  int    `json:"regular"`
	Priority int    `json:"priority"`
}

type responseData struct {
	Priority      int `json:"priority"`
	Normal        int `json:"normal"`
	Economic      int `json:"economic"`
	SuperEconomic int `json:"superEconomic"`
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

	resp, err = httpClient.Get(bcinfoAPI)
	if err != nil {
		return nil, err
	}

	superData := &sourceDataBCI{}
	err = json.NewDecoder(resp.Body).Decode(superData)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&responseData{
		Priority:      feeData.Priority,
		Normal:        feeData.Normal,
		Economic:      feeData.Economic,
		SuperEconomic: superData.Limits.Min,
	})
}
