package coinbasepro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
	"golang.org/x/time/rate"
	"net/http"
	"strconv"
	"time"
)

const (

	// Base URL's
	coinbaseproAPIURL       = "https://api.pro.coinbase.com"
	coinbaseproWebsocketURL = "wss://ws-feed.exchange.coinbase.com"

	// Base URL's for sandbox environment
	coinbaseproSandboxWebsiteURL   = "https://public.sandbox.exchange.coinbase.com"
	coinbaseproSandboxRestAPIURL   = "https://api-public.sandbox.exchange.coinbase.com"
	coinbaseproSandboxWebsocketURL = "wss://ws-feed-public.sandbox.exchange.coinbase.com"
	coinbaseproSandboxFixAPIURL    = "tcp+ssl://fix-public.sandbox.exchange.coinbase.com:4198"

	// Endpoints
	coinbaseproAccounts                = "accounts"
	coinbaseproProducts                = "products"
	coinbaseproOrderbook               = "book"
	coinbaseproTicker                  = "ticker"
	coinbaseproTrades                  = "trades"
	coinbaseproHistory                 = "candles"
	coinbaseproStats                   = "stats"
	coinbaseproCurrencies              = "currencies"
	coinbaseproLedger                  = "ledger"
	coinbaseproHolds                   = "holds"
	coinbaseproOrders                  = "orders"
	coinbaseproFills                   = "fills"
	coinbaseproTransfers               = "transfers"
	coinbaseproReports                 = "reports"
	coinbaseproTime                    = "time"
	coinbaseproFees                    = "fees"
	coinbaseproConversions             = "conversions"
	coinbaseproProfiles                = "profiles"
	coinbaseproProfilesTransfer        = "profiles/transfer"
	coinbaseproMarginTransfer          = "profiles/margin-transfer"
	coinbaseproPosition                = "position"
	coinbaseproPositionClose           = "position/close"
	coinbaseproPaymentMethod           = "payment-methods"
	coinbaseproPaymentMethodDeposit    = "deposits/payment-method"
	coinbaseproDepositCoinbase         = "deposits/coinbase-account"
	coinbaseproWithdrawalPaymentMethod = "withdrawals/payment-method"
	coinbaseproWithdrawalCoinbase      = "withdrawals/coinbase"
	coinbaseproWithdrawalCoinbaseAcct  = "withdrawals/coinbase-account"
	coinbaseproWithdrawalCrypto        = "withdrawals/crypto"
	coinbaseproWithdrawalFeeEstimate   = "withdrawals/fee-estimate"
	coinbaseproCoinbaseAccounts        = "coinbase-accounts"
	coinbaseproTrailingVolume          = "users/self/trailing-volume"
	coinbaseproExchangeLimits          = "users/self/exchange-limits"
)

type CoinbasePro struct {
	qsx.Exchange
	Conn *websocket.Conn
}

func NewCoinbasePro(config *qsx.Config) (qsx.IExchange, error) {

	t := transport{
		authKey:        config.Key,
		authPassphrase: config.Passphrase,
		authSecret:     config.Secret,
		timestamp: func() string {
			return strconv.FormatInt(time.Now().Unix(), 10)
		},
	}

	rl := rate.NewLimiter(rate.Every(time.Second), 10) // 10 requests per second

	var apiUrl string
	var wsUrl string

	if config.Sandbox {
		apiUrl = coinbaseproSandboxRestAPIURL
		wsUrl = coinbaseproSandboxWebsocketURL
	} else {
		apiUrl = coinbaseproAPIURL
		wsUrl = coinbaseproWebsocketURL
	}

	return &CoinbasePro{
		qsx.Exchange{
			Name: qsx.CoinbasePro,
			Features: &qsx.ExchangeFeatures{
				HasCrypto:    true,
				HasWebsocket: true,
				HasOptions:   false,
			},
			Auth: config.Auth,
			API: qsx.New(
				&http.Client{
					Transport:     &t,
					CheckRedirect: nil,
					Jar:           nil,
					Timeout:       0,
				},
				qsx.Options{
					ApiURL:  apiUrl,
					Verbose: false,
				},
				rl,
			),
			Websocket: &qsx.Dialer{
				URL: wsUrl,
			},
		},
		&websocket.Conn{},
	}, nil
}

type transport struct {
	authKey        string
	authPassphrase string
	authSecret     string
	timestamp      func() string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {

	var b bytes.Buffer
	if req.Body != nil {
		err := json.NewEncoder(&b).Encode(req.Body)
		if err != nil {
			return nil, fmt.Errorf("qsx coinbase: error encoding content: %w", err)
		}
	}
	timestamp := t.timestamp()
	msg := fmt.Sprintf("%s%s%s%s", timestamp, req.Method, req.URL, b.Bytes())
	signature, err := qsx.SignSHA256HMAC(msg, t.authSecret)
	if err != nil {
		return nil, fmt.Errorf("qsx coinbase: error signing content: %w", err)
	}

	r := req.Clone(req.Context())
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Accept-Charset", "UTF-8")
	r.Header.Add("User-Agent", "qsx v0.0.1")
	r.Header.Add("CB-ACCESS-KEY", t.authKey)
	r.Header.Add("CB-ACCESS-PASSPHRASE", t.authPassphrase)
	r.Header.Add("CB-ACCESS-TIMESTAMP", timestamp)
	r.Header.Add("CB-ACCESS-SIGN", signature)

	return http.DefaultTransport.RoundTrip(r)
}
