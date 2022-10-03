package tests

import (
	quantstopexchange "github.com/quantstop/quantstopterminal/pkg/exchange"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
	"testing"
)

var key = ""
var pass = ""
var secret = ""

func TestNewClient(t *testing.T) {
	config := &qsx.Config{
		Auth:    qsx.NewAuth(key, pass, secret),
		Sandbox: true,
	}

	for _, x := range qsx.SupportedExchanges {
		ex, err := quantstopexchange.NewExchange(x, config)
		if err != nil {
			t.Error(err)
		}
		t.Logf("Exchange Name: %v", ex.GetName())
	}
}
