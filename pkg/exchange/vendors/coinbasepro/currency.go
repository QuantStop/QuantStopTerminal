package coinbasepro

import (
	"context"
	"fmt"
)

type CurrencyName string

type Currency struct {
	ConvertibleTo []CurrencyName  `json:"convertible_to"`
	Details       CurrencyDetails `json:"details"`
	ID            string          `json:"id"`
	MaxPrecision  float64         `json:"max_precision"`
	Message       string          `json:"message"`
	MinSize       float64         `json:"min_size"`
	Name          string          `json:"name"`
	Status        string          `json:"status"`
}

type CurrencyDetails map[string]interface{}

// ListCurrencies retrieves the list of known Currencies. Not all Currencies may be available for trading.
func (c *CoinbasePro) ListCurrencies(ctx context.Context) ([]Currency, error) {
	var currencies []Currency
	path := fmt.Sprintf("/%s/", coinbaseproCurrencies)
	return currencies, c.API.Get(ctx, path, &currencies)
}

// GetCurrency retrieves the details of a specific Currency.
func (c *CoinbasePro) GetCurrency(ctx context.Context, currencyName CurrencyName) (Currency, error) {
	var currency Currency
	path := fmt.Sprintf("/%s/%s/", coinbaseproCurrencies, currencyName)
	return currency, c.API.Get(ctx, path, &currency)
}
