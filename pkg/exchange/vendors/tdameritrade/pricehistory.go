package tdameritrade

import (
	"context"
	"fmt"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
	"time"
)

var (
	validPeriodTypes    = []string{"day", "month", "year", "ytd"}
	validFrequencyTypes = []string{"minute", "daily", "weekly", "monthly"}
)

const (
	defaultPeriodType    = "day"
	defaultFrequencyType = "minute"
)

// PriceHistoryOptions is parsed and translated to query options in the https request
type PriceHistoryOptions struct {
	PeriodType            string `url:"periodType"`
	Period                int    `url:"period"`
	FrequencyType         string `url:"frequencyType"`
	Frequency             int    `url:"frequency"`
	EndDate               Time   `url:"endDate,omitempty"`
	StartDate             Time   `url:"startDate,omitempty"`
	NeedExtendedHoursData bool   `url:"needExtendedHoursData"`
}

type PriceHistory struct {
	Candles []struct {
		Close    float64 `json:"close"`
		Datetime int     `json:"datetime"`
		High     float64 `json:"high"`
		Low      float64 `json:"low"`
		Open     float64 `json:"open"`
		Volume   float64 `json:"volume"`
	} `json:"candles"`
	Empty  bool   `json:"empty"`
	Symbol string `json:"symbol"`
}

func (opts *PriceHistoryOptions) validate() error {
	if opts.PeriodType != "" {
		if !contains(opts.PeriodType, validPeriodTypes) {
			return fmt.Errorf("invalid periodType, must have the value of one of the following %v", validPeriodTypes)
		}
	} else {
		opts.PeriodType = defaultPeriodType
	}

	if opts.FrequencyType != "" {
		if !contains(opts.FrequencyType, validFrequencyTypes) {
			return fmt.Errorf("invalid frequencyType, must have the value of one of the following %v", validFrequencyTypes)
		}
	} else {
		opts.PeriodType = defaultFrequencyType
	}

	return nil
}

func contains(s string, lst []string) bool {
	for _, e := range lst {
		if e == s {
			return true
		}
	}
	return false
}

func (opts *PriceHistoryOptions) Params() []string {
	params := []string{fmt.Sprintf("periodType=%s", opts.PeriodType)}
	params = append(params, fmt.Sprintf("period=%v", opts.Period))
	params = append(params, fmt.Sprintf("frequencyType=%v", opts.FrequencyType))
	params = append(params, fmt.Sprintf("frequency=%v", opts.Frequency))
	params = append(params, fmt.Sprintf("endDate=%v", opts.EndDate))
	params = append(params, fmt.Sprintf("startDate=%v", opts.StartDate))
	params = append(params, fmt.Sprintf("needExtendedHoursData=%v", &opts.NeedExtendedHoursData))

	if !opts.EndDate.Time().IsZero() {
		end := opts.EndDate.Time().Format(time.RFC3339Nano)
		params = append(params, fmt.Sprintf("endDate=%s", end))
	}
	if !opts.StartDate.Time().IsZero() {
		start := opts.StartDate.Time().Format(time.RFC3339Nano)
		params = append(params, fmt.Sprintf("startDate=%s", start))
	}

	return params
}

// PriceHistory get the price history for a symbol
// TDAmeritrade API Docs: https://developer.tdameritrade.com/price-history/apis/get/marketdata/%7Bsymbol%7D/pricehistory
func (c *TDAmeritrade) PriceHistory(ctx context.Context, symbol string, opts PriceHistoryOptions) (PriceHistory, error) {
	var history PriceHistory

	if err := opts.validate(); err != nil {
		return history, err
	}

	path := fmt.Sprintf("/%s/%s/%s/%s", marketdata, symbol, priceHistory, qsx.Query(opts.Params()))
	if err := c.API.Get(ctx, path, &history); err != nil {
		return history, err
	}
	return history, nil
}
