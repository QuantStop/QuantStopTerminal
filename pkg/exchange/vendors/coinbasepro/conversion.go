package coinbasepro

import (
	"context"
	"encoding/json"
	"fmt"
)

type StablecoinConversionSpec struct {
	From   CurrencyName `json:"from"`
	To     CurrencyName `json:"to"`
	Amount float64      `json:"amount,string"`
}

func (s *StablecoinConversionSpec) UnmarshalJSON(b []byte) error {
	type Alias StablecoinConversionSpec
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	return json.Unmarshal(b, &aux)
}

type StablecoinConversion struct {
	Amount        float64      `json:"amount"`
	From          CurrencyName `json:"from"`
	FromAccountID string       `json:"from_account_id"`
	ID            string       `json:"id"`
	To            CurrencyName `json:"to"`
	ToAccountID   string       `json:"to_account_id"`
}

// CreateStablecoinConversion creates a conversion from a crypto Currency a stablecoin Currency.
func (c *CoinbasePro) CreateStablecoinConversion(ctx context.Context, stablecoinConversionSpec StablecoinConversionSpec) (StablecoinConversion, error) {
	var result StablecoinConversion
	path := fmt.Sprintf("/%s/", coinbaseproConversions)
	return result, c.API.Post(ctx, path, stablecoinConversionSpec, &result)
}
