package handlers

/*import (
	"context"
	"github.com/quantstop/quantstopterminal/pkg/quantstopexchange/qsx"
	"github.com/quantstop/quantstopterminal/internal"
	"github.com/quantstop/quantstopterminal/internal/database/models"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/internal/webserver/errors"
	"github.com/quantstop/quantstopterminal/internal/webserver/router"
	"github.com/quantstop/quantstopterminal/internal/webserver/write"
	"net/http"
)

type getCandleResponse struct {
	Type          string       `json:"type"`
	HistoricRates []qsx.Candle `json:"candles"`
}

// GetCandles
// Historic rates for a product.
// Rates are returned in grouped buckets.
// Candle schema is of the form [timestamp, price_low, price_high, price_open, price_close]
// Request: GET "/api/exchanges/([^/]+)/products/([^/]+)/candles"
// Params:
// - granularity (string, required) {60, 300, 900, 3600, 21600, 86400}
// - start (string, optional)
// - end (string, optional)
// Example: GET "/api/exchanges/coinbase/products/BTC-USD/candles?granularity=5?"
func GetCandles(bot internal.IEngine, user *models.User, w http.ResponseWriter, r *http.Request) http.HandlerFunc {

	exchangeName := router.GetField(r, 0)
	productName := router.GetField(r, 1)
	granularity := r.URL.Query().Get("granularity")

	historicalCandles, err := bot.GetExchange(exchangeName).GetHistoricalCandles(context.TODO(), productName, granularity)
	if err != nil {
		log.Debugln(log.Webserver, "error getting candles: %v", err)
		return write.Error(errors.InternalError)
	}
	resp := getCandleResponse{
		Type:          "getProductCandles",
		HistoricRates: historicalCandles,
	}
	return write.JSON(resp)

}

type getProductResponse struct {
	Type     string        `json:"type"`
	Products []qsx.Product `json:"products"`
}

// GetProducts
// Supported Currency Pairs for an Exchange.
// Request: GET "/api/exchanges/([^/]+)/products"
// Params:
// Example: GET "/api/exchanges/coinbase/products
func GetProducts(bot internal.IEngine, user *models.User, w http.ResponseWriter, r *http.Request) http.HandlerFunc {

	exchangeName := router.GetField(r, 0)
	log.Debugf(log.Webserver, "GetProducts: %v", exchangeName)

	if exchangeName == "" {
		return write.JSON(errors.InternalError)
	}

	exchange := bot.GetExchange(exchangeName)

	products, err := exchange.ListProducts(context.TODO())
	if err != nil {
		log.Debugf(log.Webserver, "error getting products: %v", err)
	}
	resp := getProductResponse{
		Type:     "getProducts",
		Products: products,
	}
	return write.JSON(resp)

}*/
