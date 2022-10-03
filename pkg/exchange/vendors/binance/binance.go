package binance

import (
	"context"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

const (
	apiURL         = "https://api.binance.com"
	spotAPIURL     = "https://sapi.binance.com"
	cfuturesAPIURL = "https://dapi.binance.com"
	ufuturesAPIURL = "https://fapi.binance.com"

	// Public endpoints
	exchangeInfo      = "/api/v3/exchangeInfo"
	orderBookDepth    = "/api/v3/depth"
	recentTrades      = "/api/v3/trades"
	aggregatedTrades  = "/api/v3/aggTrades"
	candleStick       = "/api/v3/klines"
	averagePrice      = "/api/v3/avgPrice"
	priceChange       = "/api/v3/ticker/24hr"
	symbolPrice       = "/api/v3/ticker/price"
	bestPrice         = "/api/v3/ticker/bookTicker"
	userAccountStream = "/api/v3/userDataStream"
	perpExchangeInfo  = "/fapi/v1/exchangeInfo"
	historicalTrades  = "/api/v3/historicalTrades"

	// Authenticated endpoints
	newOrderTest      = "/api/v3/order/test"
	orderEndpoint     = "/api/v3/order"
	openOrders        = "/api/v3/openOrders"
	allOrders         = "/api/v3/allOrders"
	accountInfo       = "/api/v3/account"
	marginAccountInfo = "/sapi/v1/margin/account"

	// Withdraw API endpoints
	accountStatus                          = "/wapi/v3/accountStatus.html"
	systemStatus                           = "/wapi/v3/systemStatus.html"
	dustLog                                = "/wapi/v3/userAssetDribbletLog.html"
	tradeFee                               = "/wapi/v3/tradeFee.html"
	assetDetail                            = "/wapi/v3/assetDetail.html"
	undocumentedInterestHistory            = "/gateway-api/v1/public/isolated-margin/pair/vip-level"
	undocumentedCrossMarginInterestHistory = "/gateway-api/v1/friendly/margin/vip/spec/list-all"

	// Wallet endpoints
	allCoinsInfo     = "/sapi/v1/capital/config/getall"
	withdrawEndpoint = "/sapi/v1/capital/withdraw/apply"
	depositHistory   = "/sapi/v1/capital/deposit/hisrec"
	withdrawHistory  = "/sapi/v1/capital/withdraw/history"
	depositAddress   = "/sapi/v1/capital/deposit/address"

	defaultRecvWindow     = 5 * time.Second
	binanceSAPITimeLayout = "2006-01-02 15:04:05"
)

// OAuth example
var (
	authConfig = oauth2.Config{
		ClientID:     "XXXX-XXXX-XXXX-XXXX",
		ClientSecret: "YYYY-YYYY-YYYY-YYYY",
		RedirectURL:  "https://api.our-service.com/oauth/callback",
		Scopes:       []string{"all"},
		Endpoint: oauth2.Endpoint{
			AuthStyle: oauth2.AuthStyleInParams,
			AuthURL:   "https://api.serbvice-b.com/oauth/authorize",
			TokenURL:  "https://api.serbvice-b.com/oauth/access_token",
		},
	}
)

type Binance struct {
	qsx.Exchange
}

func NewBinance(config *qsx.Config) (qsx.IExchange, error) {

	// OAuth example
	httpClient := authConfig.Client(
		context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: &transport{}}),
		config.Token,
	)

	rl := rate.NewLimiter(rate.Every(time.Second), 10) // 10 requests per second

	api := qsx.New(
		httpClient,
		qsx.Options{
			ApiURL:  apiURL,
			Verbose: false,
		},
		rl,
	)

	return &Binance{
		qsx.Exchange{
			Name: qsx.Binance,
			Auth: config.Auth,
			API:  api,
		},
	}, nil
}

// OAuth example
type transport struct{}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")
	return http.DefaultTransport.RoundTrip(r)
}
