package tests

import (
	"context"
	"fmt"
	quantstopexchange "github.com/quantstop/quantstopterminal/pkg/exchange"

	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
	"github.com/quantstop/quantstopterminal/pkg/exchange/vendors/coinbasepro"
	"golang.org/x/sync/errgroup"
	"sync"
	"testing"
)

var cbpKey = ""
var cbpPass = ""
var cbpSecret = ""

var CoinbaseproClient qsx.IExchange
var err error

func TestCoinbaseClient(t *testing.T) {
	config := &qsx.Config{
		Auth:    qsx.NewAuth(cbpKey, cbpPass, cbpSecret),
		Sandbox: false,
	}
	CoinbaseproClient, err = quantstopexchange.NewExchange(qsx.CoinbasePro, config)
	if err != nil {
		t.Error(err)
	}
}

func TestCoinbaseCandles(t *testing.T) {
	TestCoinbaseClient(t)
	candles, err := CoinbaseproClient.GetHistoricalCandles(context.TODO(), "BTC-USD", "1m")
	if err != nil {
		t.Error(err)
	}
	for _, candle := range candles {
		t.Logf("Candle Time: %v | Open: %v | High: %v | Low: %v | Close: %v | Volume: %v", candle.Time, candle.Open, candle.High, candle.Low, candle.Close, candle.Volume)
	}
}

func TestCoinbaseListProducts(t *testing.T) {
	TestCoinbaseClient(t)
	products, err := CoinbaseproClient.ListProducts(context.TODO())
	if err != nil {
		t.Error(err)
	}
	for _, product := range products {
		t.Logf("Product: %v", product.ID)
	}
}

func TestCoinbaseFeed(t *testing.T) {
	TestCoinbaseClient(t)
	feed := coinbasepro.NewFeed()

	ctx := context.TODO()
	wg, ctx := errgroup.WithContext(ctx)

	s := &sync.WaitGroup{}
	shutdown := make(chan struct{})

	// start api client feed
	book, err := CoinbaseproClient.WatchFeed(shutdown, s, "BTC-USD", feed)
	if err != nil {
		t.Error(err)
	}

	// Loop on Heartbeat channel
	wg.Go(func() error {
		for message := range feed.Heartbeat {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				out := fmt.Sprintf("%s | %s | %s | %v | %v", message.Type, message.Time.String(), message.ProductId, message.Sequence, message.LastTradeId)
				fmt.Println(out)

				bk := fmt.Sprintf("Best Bid: %v | Best Ask: %v ", book.GetBestBid(), book.GetBestOffer())
				fmt.Println(bk)
			}
		}
		return nil
	})

	// Loop on L2Channel channel
	wg.Go(func() error {
		for message := range feed.Level2 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				out := fmt.Sprintf("%s | %s | %v | %v", message.Type, message.Time.String(), message.ProductId, message.Changes)
				fmt.Println(out)
			}
		}
		return nil
	})

	// Loop on L2ChannelSnapshot channel
	wg.Go(func() error {
		for message := range feed.Level2Snap {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				out := fmt.Sprintf("%s | %s", message.Type, message.ProductId)
				fmt.Println(out)
			}
		}
		return nil
	})

	// Loop on Matches channel
	wg.Go(func() error {
		for message := range feed.Matches {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				out := fmt.Sprintf("%s | %s | %s | %v", message.Type, message.Time.String(), message.ProductId, message.Price)
				fmt.Println(out)
			}
		}
		return nil
	})

	_ = wg.Wait()

}
