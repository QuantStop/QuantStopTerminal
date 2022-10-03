package coinbasepro

import (
	"context"
	"fmt"
)

// Limits provide payment method transfer limits, as well as buy/sell limits per currency.
type Limits struct {
	LimitCurrency  CurrencyName             `json:"limit_currency"`
	TransferLimits map[string]CurrencyLimit `json:"transfer_limits"`
}

type CurrencyLimit map[CurrencyName]Limit

// TODO: haven't ever seen PeriodInDays
type Limit struct {
	Max          float64 `json:"max"`
	Remaining    float64 `json:"remaining"`
	PeriodInDays int     `json:"period_in_days"`
}

// GetLimits retrieves the payment method transfer limits and per currency buy/sell limits for the current Profile.
func (c *CoinbasePro) GetLimits(ctx context.Context) (Limits, error) {
	var limits Limits
	path := fmt.Sprintf("/%s/", coinbaseproExchangeLimits)
	return limits, c.API.Get(ctx, path, &limits)
}
