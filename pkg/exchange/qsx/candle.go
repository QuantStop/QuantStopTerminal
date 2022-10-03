package qsx

import "time"

// Candle is the main type to represent candlestick data
type Candle struct {
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Volume float64   `json:"volume"`
	Time   time.Time `json:"time"`
}
