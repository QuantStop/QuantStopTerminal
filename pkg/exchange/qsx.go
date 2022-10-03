package exchange

import (
	"errors"
	"fmt"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
	"github.com/quantstop/quantstopterminal/pkg/exchange/vendors/binance"
	"github.com/quantstop/quantstopterminal/pkg/exchange/vendors/coinbasepro"
	"github.com/quantstop/quantstopterminal/pkg/exchange/vendors/tdameritrade"
	"github.com/quantstop/quantstopterminal/pkg/exchange/vendors/yfinance"
)

// NewExchange creates an exchange connection and returns a struct that implements the IExchange interface
func NewExchange(name qsx.Name, config *qsx.Config) (qsx.IExchange, error) {

	found := false
	for _, x := range qsx.SupportedExchanges {
		if x == name {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New(fmt.Sprintf("qsx error: failed to create exchange, '%s' is not supported", name))
	}

	switch name {
	case qsx.CoinbasePro:
		c, err := coinbasepro.NewCoinbasePro(config)
		if err != nil {
			return nil, err
		}
		return c, nil

	case qsx.Binance:
		b, err := binance.NewBinance(config)
		if err != nil {
			return nil, err
		}
		return b, nil

	case qsx.YFinance:
		b, err := yfinance.NewYFinance(config)
		if err != nil {
			return nil, err
		}
		return b, nil

	case qsx.TDAmeritrade:
		b, err := tdameritrade.NewTDAmeritrade(config)
		if err != nil {
			return nil, err
		}
		return b, nil

	default:
		return nil, errors.New("qsx error: failed to create exchange, unexpected error")
	}

}

// GetSupportedExchanges returns a list of all the supported exchanges
func GetSupportedExchanges() []qsx.Name {
	return qsx.SupportedExchanges
}
