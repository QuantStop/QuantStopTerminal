package tdameritrade

import (
	"context"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx/orderbook"
	"sync"
	"time"
)

func (c *TDAmeritrade) GetHistoricalCandles(ctx context.Context, productID string, granularity string) ([]qsx.Candle, error) {
	var candles []qsx.Candle

	var p int
	var pt string
	var f int
	var ft string

	switch granularity {
	case "1m":
		p = 1
		pt = "day"
		f = 1
		ft = "minute"
	case "5m":
		p = 5
		pt = "day"
		f = 5
		ft = "minute"
	case "15m":
		p = 5
		pt = "day"
		f = 15
		ft = "minute"
	case "1h":
		p = 10
		pt = "day"
		f = 1
		ft = "daily"
	case "6h":
		p = 10
		pt = "day"
		f = 6
		ft = "daily"
	case "1d":
		p = 1
		pt = "year"
		f = 1
		ft = "daily"
	default:
		p = 1
		pt = "day"
		f = 1
		ft = "minute"

	}

	tdCandles, err := c.PriceHistory(ctx, productID, PriceHistoryOptions{
		PeriodType:            pt,
		Period:                p,
		FrequencyType:         ft,
		Frequency:             f,
		EndDate:               Time{},
		StartDate:             Time{},
		NeedExtendedHoursData: false,
	})
	if err != nil {
		return nil, err
	}

	for _, tdCandle := range tdCandles.Candles {
		candles = append(candles, qsx.Candle{
			Open:   tdCandle.Open,
			High:   tdCandle.High,
			Low:    tdCandle.Low,
			Close:  tdCandle.Close,
			Volume: tdCandle.Volume,
			Time:   time.UnixMilli(int64(tdCandle.Datetime)), //todo: ???
		})
	}

	return candles, nil
}

func (c *TDAmeritrade) WatchFeed(shutdown chan struct{}, wg *sync.WaitGroup, product string, feed interface{}) (*orderbook.Orderbook, error) {
	return nil, nil
}

func (c *TDAmeritrade) ListProducts(ctx context.Context) ([]qsx.Product, error) {
	return nil, nil
}
