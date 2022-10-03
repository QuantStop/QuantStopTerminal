package coinbasepro

import (
	"context"
	"fmt"
)

// Fees describes the current maker & taker fee rates, as well as the 30-day trailing volume.
// Quoted rates are subject to change.
// Note: the docs (https://docs.pro.coinbase.com/#fees) are wrong; the response is an object, not an array
type Fees struct {
	MakerFeeRate float64 `json:"maker_fee_rate"`
	TakerFeeRate float64 `json:"taker_fee_rate"`
	USDVolume    float64 `json:"usd_volume"`
}

// GetFees returns current maker & taker fee rates, as well as the 30-day trailing volume. GetFees is plural, but returns
// a single object. Perhaps there is a better name.
func (c *CoinbasePro) GetFees(ctx context.Context) (Fees, error) {
	var fees Fees
	path := fmt.Sprintf("/%s/", coinbaseproFees)
	return fees, c.API.Get(ctx, path, &fees)
}
