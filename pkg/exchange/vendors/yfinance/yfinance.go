package yfinance

import (
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

type YFinance struct {
	qsx.Exchange
}

func NewYFinance(config *qsx.Config) (qsx.IExchange, error) {

	rl := rate.NewLimiter(rate.Every(time.Second), 10) // 10 requests per second

	api := qsx.New(
		&http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		},
		qsx.Options{
			ApiURL:  "",
			Verbose: false,
		},
		rl,
	)

	return &YFinance{
		qsx.Exchange{
			Name: qsx.YFinance,
			Features: &qsx.ExchangeFeatures{
				HasCrypto:    false,
				HasWebsocket: true,
				HasOptions:   true,
			},
			Auth: config.Auth,
			API:  api,
		},
	}, nil
}
