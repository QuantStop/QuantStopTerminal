package binance

import (
	"context"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx/orderbook"
	"sync"
)

func (b *Binance) GetHistoricalCandles(ctx context.Context, productID string, granularity string) ([]qsx.Candle, error) {
	return nil, nil
}

func (b *Binance) WatchFeed(shutdown chan struct{}, wg *sync.WaitGroup, product string, feed interface{}) (*orderbook.Orderbook, error) {
	return nil, nil
}

func (b *Binance) ListProducts(ctx context.Context) ([]qsx.Product, error) {
	return nil, nil
}
