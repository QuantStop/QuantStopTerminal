package coinbasepro

import (
	"context"
	"errors"
	"fmt"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx/orderbook"
	"sync"
	"time"
)

// This file holds all methods that implement the main IExchange interface.
// This is served as an adapter to the Coinbase client, which adheres strictly to the Coinbase API.
// Here is where the data is formatted into the common type defined by IExchange.

func (c *CoinbasePro) GetHistoricalCandles(ctx context.Context, productID string, granularity string) ([]qsx.Candle, error) {
	var candles []qsx.Candle

	g := TimePeriod1Minute
	switch granularity {
	case "1m":
		g = TimePeriod1Minute
	case "5m":
		g = TimePeriod5Minutes
	case "15m":
		g = TimePeriod15Minutes
	case "1h":
		g = TimePeriod1Hour
	case "6h":
		g = TimePeriod6Hours
	case "1d":
		g = TimePeriod1Day
	default:
		g = TimePeriod1Minute

	}

	coinbaseCandles, err := c.GetHistoricRates(ctx, productID, HistoricRateFilter{
		Granularity: g,
		End:         Time{},
		Start:       Time{},
	})
	if err != nil {
		return nil, err
	}

	for _, cbCandle := range coinbaseCandles.Candles {
		candles = append(candles, qsx.Candle{
			Close:  cbCandle.Close,
			High:   cbCandle.High,
			Low:    cbCandle.Low,
			Open:   cbCandle.Open,
			Time:   time.Time(cbCandle.Time),
			Volume: cbCandle.Volume,
		})
	}

	return candles, nil
}

func (c *CoinbasePro) WatchFeed(shutdown chan struct{}, wg *sync.WaitGroup, product string, feed interface{}) (*orderbook.Orderbook, error) {

	// check if exchange has support for websocket streaming
	if !c.Features.HasWebsocket {
		return nil, errors.New(fmt.Sprintf("qsx error: exchange, '%s' does not have websocket support", c.Name))
	}

	// create a new subscription request
	prods := []ProductID{ProductID(product)}
	channelNames := []ChannelName{
		ChannelNameHeartbeat,
		ChannelNameLevel2,
	}
	channels := []Channel{
		{
			Name:       ChannelNameMatches,
			ProductIDs: []ProductID{ProductID(product)},
		},
	}
	subReq := NewSubscriptionRequest(prods, channelNames, channels)

	return c.Watch(shutdown, wg, subReq, feed.(*Feed))
}

func (c *CoinbasePro) ListProducts(ctx context.Context) ([]qsx.Product, error) {
	products, err := c.ListCoinbaseProducts(ctx)
	if err != nil {
		return nil, err
	}
	var returnArr []qsx.Product
	for _, product := range products {
		returnArr = append(returnArr, qsx.Product{
			ID:             product.ID,
			BaseCurrency:   string(product.BaseCurrency),
			QuoteCurrency:  string(product.QuoteCurrency),
			BaseMinSize:    product.BaseMinSize,
			BaseMaxSize:    product.BaseMaxSize,
			QuoteIncrement: product.QuoteIncrement,
			BaseIncrement:  product.BaseIncrement,
			DisplayName:    product.DisplayName,
			MinMarketFunds: product.MinMarketFunds,
			MaxMarketFunds: product.MaxMarketFunds,
			MarginEnabled:  product.MarginEnabled,
			PostOnly:       product.PostOnly,
			LimitOnly:      product.LimitOnly,
			CancelOnly:     product.CancelOnly,
			Status:         string(product.Status),
			StatusMessage:  product.StatusMessage,
			AuctionMode:    product.AuctionMode,
		})
	}

	return returnArr, nil
}
