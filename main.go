package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
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

type Candles struct {
	Start  string  `json:"start"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

func (c *Candles) UnmarshalJSON(data []byte) error {
	// Define an intermediate struct to hold the raw JSON data
	type RawCandles struct {
		Start  string `json:"start"`
		Low    string `json:"low"`
		High   string `json:"high"`
		Open   string `json:"open"`
		Close  string `json:"close"`
		Volume string `json:"volume"`
	}

	// Unmarshal the raw JSON data into the intermediate struct
	var raw RawCandles
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	// Parse the string values into float64 values
	c.Start = raw.Start
	c.Low, err = strconv.ParseFloat(raw.Low, 64)
	if err != nil {
		return err
	}
	c.High, err = strconv.ParseFloat(raw.High, 64)
	if err != nil {
		return err
	}
	c.Open, err = strconv.ParseFloat(raw.Open, 64)
	if err != nil {
		return err
	}
	c.Close, err = strconv.ParseFloat(raw.Close, 64)
	if err != nil {
		return err
	}
	c.Volume, err = strconv.ParseFloat(raw.Volume, 64)
	if err != nil {
		return err
	}

	return nil
}


func coinbaseRequest(productID string, startTime int64, endTime int64, granularity Granularity) ([]Candles, error) {
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
	var candles []Candles

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

func main() {
	// Create a new logger with log rotation
	logger := &lumberjack.Logger{
		Filename:   "/var/log/coinbase-etl.log",
		MaxSize:    100, // megabytes
		MaxAge:     7,   // days
		LocalTime:  true,
		Compress:   true,
		MaxBackups: 5,
	}

    // Configure logrus to write to the log file
    log := logrus.New()
    log.SetFormatter(&logrus.JSONFormatter{})
    log.SetOutput(logger)

    // Log a message using logrus
    log.Info("Hello from Go!")

	c := cron.New()

    // Define a cron job that runs every minute
    c.AddFunc("* * * * *", func() {
        log.Info("Running cron job 1...")
    })

    // Define a cron job that runs every hour
    c.AddFunc("0 * * * *", func() {
        log.Info("Running cron job 2...")
    })

    // Start the cron scheduler
    c.Start()

}
