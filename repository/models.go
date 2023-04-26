package repository

import (
	"encoding/json"
	"strconv"

	"gorm.io/gorm"
)

type Candles struct {
	gorm.Model
	Start  string  `json:"start"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

// Define a custom UnmarshalJSON method on the Candles struct
// when pulling from coinbase API
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
