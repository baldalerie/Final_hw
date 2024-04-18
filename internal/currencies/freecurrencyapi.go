package currencies

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var apikey = "fca_live_splQmmKr8wXgu8kP7PfeLqUna4WsWDkBVJZqy253"
var client = httpClient()

const BaseUrl = "https://api.freecurrencyapi.com/v1/"

type CurrencyResponse struct {
	Data map[string]float64 `json:"data"`
}

func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

func apiCall(endpoint string, params ...map[string]string) []byte {
	if len(apikey) == 0 {
		log.Fatalf("No API key provided!")
	}

	jsonParams, err := json.Marshal(params)

	req, err := http.NewRequest("GET", BaseUrl+endpoint, bytes.NewBuffer(jsonParams))
	if err != nil {
		log.Fatalf("Error Occurred. %+v", err)
	}

	req.Header.Set("apikey", apikey)

	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request to API endpoint. %+v", err)
	}

	// Close the connection to reuse it
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Couldn't parse response body. %+v", err)
	}

	return body
}

func apiCallData(endpoint string, params ...map[string]string) (*CurrencyResponse, error) {
	if len(apikey) == 0 {
		log.Fatalf("No API key provided!")
	}

	req, err := http.NewRequest("GET", BaseUrl+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("apikey", apikey)
	//fmt.Println("url Currency: ", req)
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to API endpoint: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", response.Status)
	}

	var currencyResp CurrencyResponse
	if err := json.NewDecoder(response.Body).Decode(&currencyResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &currencyResp, nil
}

func Status() []byte {
	return apiCall("status")
}

func Currencies(params map[string]string) []byte {
	return apiCall("currencies", params)
}

func Latest(params map[string]string) []byte {
	return apiCall("latest", params)
}

func LatestData(params map[string]string) (*CurrencyResponse, error) {
	return apiCallData("latest", params)
}

func Historical(params map[string]string) []byte {
	return apiCall("historical", params)
}
