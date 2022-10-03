package qsx

// Name is the unique name of the exchange
type Name string

const (
	CoinbasePro  Name = "coinbasepro"
	Binance      Name = "binance"
	YFinance     Name = "yfinance"
	TDAmeritrade Name = "tdameritrade"
)

// SupportedExchanges is a list of all supported exchange connections
var SupportedExchanges = []Name{
	CoinbasePro,
	//TDAmeritrade,
}
