package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
}

func GetUAHExchangeRate() (float64, error) {
	url := "https://api.exchangerate-api.com/v4/latest/USD?apikey=%s"

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Failed to get exchange rate:", err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body:", err)
		return 0, err
	}

	var data ExchangeRateResponse
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println("Failed to parse JSON:", err)
		return 0, err
	}

	uahRate, exists := data.Rates["UAH"]
	if !exists {
		return 0, fmt.Errorf("UAH exchange rate not found in response")
	}

	return uahRate, nil
}
