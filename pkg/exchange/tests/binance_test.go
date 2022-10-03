package tests

/*import (
	"context"
	quantstopexchange "github.com/quantstop/quantstopterminal/pkg/exchange"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
	"testing"
)

var binanceKey = ""
var binancePass = ""
var binanceSecret = ""

var BinanceClient qsx.IExchange

func TestBinanceClient(t *testing.T) {
	config := &qsx.Config{
		Auth:    qsx.NewAuth(binanceKey, binancePass, binanceSecret),
		Sandbox: true,
	}
	BinanceClient, err = quantstopexchange.NewExchange(qsx.Binance, config)
	if err != nil {
		t.Error(err)
	}
}

func TestBinanceCandles(t *testing.T) {
	TestBinanceClient(t)
	candles, err := BinanceClient.GetHistoricalCandles(context.TODO(), "BTC-USD", "1m")
	if err != nil {
		t.Error(err)
	}
	for _, candle := range candles {
		t.Logf("Candle Time: %v | Open: %v | High: %v | Low: %v | Close: %v | Volume: %v", candle.Time, candle.Open, candle.High, candle.Low, candle.Close, candle.Volume)
	}
}*/
