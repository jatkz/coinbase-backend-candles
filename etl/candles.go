package etl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	repo "coinbase-etl/repository"
)

type Granularity struct {
	Granularity string
}

func (g Granularity) String() string {
	return g.Granularity
}

var (
	GranularityOneMinute  = Granularity{"ONE_MINUTE"}
	GranularityFiveMinute = Granularity{"FIVE_MINUTE"}
	GranularityFifteenMinute = Granularity{"FIFTEEN_MINUTE"}
	GranularityThirtyMinute = Granularity{"THIRTY_MINUTE"}
	GranularityOneHour =  Granularity{"ONE_HOUR"}
	GranularityTwoHour =  Granularity{"TWO_HOUR"}
	GranularitySixHour = Granularity{"SIX_HOUR"}
	GranularityOneDay = Granularity{"ONE_DAY"}
)

// Generate the URL for the candles endpoint
func buildCandlesURL(productID string, startTime int64, endTime int64, granularity Granularity) (string, error) {
	// Define the base URL
	baseURL := "https://api.coinbase.com/api/v3/brokerage/products/:product_id/candles"

	// Define the parameters
	params := url.Values{}
	params.Set("start", fmt.Sprintf("%d", startTime))
	params.Set("end", fmt.Sprintf("%d", endTime))
	params.Set("granularity", granularity.String())

	// Replace the :product_id placeholder with the actual product ID
	url := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	url = strings.ReplaceAll(url, ":product_id", productID)

	return url, nil
}

// Simple method for setting the endTime to the current time
func getCurrentUnixTimestamp() int64 {
	return time.Now().Unix()
}

// Simple method for setting the start date
func getUnixTimestampDaysAgo(daysAgo int) int64 {
	// Get the current time
	now := time.Now()

	// Calculate the duration for the specified number of days ago
	duration := time.Duration(daysAgo) * 24 * time.Hour

	// Subtract the duration from the current time to get the timestamp for the specified number of days ago
	return now.Add(-duration).Unix()
}

func CoinbaseRequest(productID string, startTime int64, endTime int64, granularity Granularity) ([]repo.Candles, error) {
	url, err := buildCandlesURL(productID, startTime, endTime, granularity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}	
	
	method := "GET"
  
	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)
  
	if err != nil {
	  fmt.Println(err)
	  return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
  
	res, err := client.Do(req)
	if err != nil {
	  fmt.Println(err)
	  return nil, err
	}
	defer res.Body.Close()
  
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
	  fmt.Println(err)
	  return nil, err
	}
	fmt.Println(string(body))

	// Declare a slice of Person objects to hold the parsed data
	var candles []repo.Candles

	// Use json.Unmarshal to parse the JSON data into the slice of Person objects
	err = json.Unmarshal(body, &candles)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, err
	}

	// Print the parsed data to the console
	fmt.Println("Parsed JSON data:")
	for _, candle := range candles {
		fmt.Println(candle.Start, candle.Low, candle.High, candle.Open, candle.Close, candle.Volume)
	}
	
	return candles, nil
}
