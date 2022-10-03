package tests

/*import (
	"context"
	quantstopexchange "github.com/quantstop/quantstopterminal/pkg/exchange"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
	"testing"
)

var tdaKey = ""
var tdaPass = ""
var tdaSecret = ""

var TDAClient qsx.IExchange
var tdaErr error

func TestTDAClient(t *testing.T) {
	config := &qsx.Config{
		Auth:    qsx.NewAuth(tdaKey, tdaPass, tdaSecret),
		Sandbox: false,
	}
	TDAClient, tdaErr = quantstopexchange.NewExchange(qsx.TDAmeritrade, config)
	if tdaErr != nil {
		t.Error(tdaErr)
	}
}

func TestTDACandles(t *testing.T) {
	TestTDAClient(t)

	candles, err := TDAClient.GetHistoricalCandles(context.TODO(), "GME", "1m")
	if err != nil {
		t.Error(err)
	}

	for _, candle := range candles {
		t.Logf("Candle Time: %v | Open: %v | High: %v | Low: %v | Close: %v | Volume: %v", candle.Time, candle.Open, candle.High, candle.Low, candle.Close, candle.Volume)
	}
}
*/
