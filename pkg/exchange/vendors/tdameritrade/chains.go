package tdameritrade

import (
	"context"
	"fmt"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
	"time"
)

// TDAmeritrade Options Chains API docs: https://developer.tdameritrade.com/option-chains/apis/get/marketdata/chains

type Underlying struct {
	Symbol            string  `json:"symbol"`
	Description       string  `json:"description"`
	Change            float64 `json:"change"`
	PercentChange     float64 `json:"percentChange"`
	Close             float64 `json:"close"`
	QuoteTime         int     `json:"quoteTime"`
	TradeTime         int     `json:"tradeTime"`
	Bid               float64 `json:"bid"`
	Ask               float64 `json:"ask"`
	Last              float64 `json:"last"`
	Mark              float64 `json:"mark"`
	MarkChange        float64 `json:"markChange"`
	MarkPercentChange float64 `json:"markPercentChange"`
	BidSize           int     `json:"bidSize"`
	AskSize           int     `json:"askSize"`
	HighPrice         float64 `json:"highPrice"`
	LowPrice          float64 `json:"lowPrice"`
	OpenPrice         float64 `json:"openPrice"`
	TotalVolume       int     `json:"totalVolume"`
	ExchangeName      string  `json:"exchangeName"`
	FiftyTwoWeekHigh  float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow   float64 `json:"fiftyTwoWeekLow"`
	Delayed           bool    `json:"delayed"`
}

type ExpDateOption struct {
	PutCall                string  `json:"putCall"`
	Symbol                 string  `json:"symbol"`
	Description            string  `json:"description"`
	ExchangeName           string  `json:"exchangeName"`
	Bid                    float64 `json:"bid"`
	Ask                    float64 `json:"ask"`
	Last                   float64 `json:"last"`
	Mark                   float64 `json:"mark"`
	BidSize                int     `json:"bidSize"`
	AskSize                int     `json:"askSize"`
	BidAskSize             string  `json:"bidAskSize"`
	LastSize               float64 `json:"lastSize"`
	HighPrice              float64 `json:"highPrice"`
	LowPrice               float64 `json:"lowPrice"`
	OpenPrice              float64 `json:"openPrice"`
	ClosePrice             float64 `json:"closePrice"`
	TotalVolume            int     `json:"totalVolume"`
	TradeDate              string  `json:"tradeDate"`
	TradeTimeInLong        int     `json:"tradeTimeInLong"`
	QuoteTimeInLong        int     `json:"quoteTimeInLong"`
	NetChange              float64 `json:"netChange"`
	Volatility             float64 `json:"volatility"`
	Delta                  float64 `json:"delta"`
	Gamma                  float64 `json:"gamma"`
	Theta                  float64 `json:"theta"`
	Vega                   float64 `json:"vega"`
	Rho                    float64 `json:"rho"`
	OpenInterest           int     `json:"openInterest"`
	TimeValue              float64 `json:"timeValue"`
	TheoreticalOptionValue float64 `json:"theoreticalOptionValue"`
	TheoreticalVolatility  float64 `json:"theoreticalVolatility"`
	OptionDeliverablesList string  `json:"optionDeliverablesList"`
	StrikePrice            float64 `json:"strikePrice"`
	ExpirationDate         int     `json:"expirationDate"`
	DaysToExpiration       int     `json:"daysToExpiration"`
	ExpirationType         string  `json:"expirationType"`
	LastTradingDate        int     `json:"lastTradingDay"`
	Multiplier             float64 `json:"multiplier"`
	SettlementType         string  `json:"settlementType"`
	DeliverableNote        string  `json:"deliverableNote"`
	IsIndexOption          bool    `json:"isIndexOption"`
	PercentChange          float64 `json:"percentChange"`
	MarkChange             float64 `json:"markChange"`
	MarkPercentChange      float64 `json:"markPercentChange"`
	InTheMoney             bool    `json:"inTheMoney"`
	Mini                   bool    `json:"mini"`
	NonStandard            bool    `json:"nonStandard"`
}

// ExpDateMap
// the first string is the exp date.
// the second string is the strike price.
type ExpDateMap map[string]map[string][]ExpDateOption

type Chains struct {
	Symbol            string     `json:"symbol"`
	Status            string     `json:"status"`
	Underlying        Underlying `json:"underlying"`
	Strategy          string     `json:"strategy"`
	Interval          float64    `json:"interval"`
	IsDelayed         bool       `json:"isDelayed"`
	IsIndex           bool       `json:"isIndex"`
	InterestRate      float64    `json:"interestRate"`
	UnderlyingPrice   float64    `json:"underlyingPrice"`
	Volatility        float64    `json:"volatility"`
	DaysToExpiration  float64    `json:"daysToExpiration"`
	NumberOfContracts int        `json:"numberOfContracts"`
	CallExpDateMap    ExpDateMap `json:"callExpDateMap"`
	PutExpDateMap     ExpDateMap `json:"putExpDateMap"`
}

type Strategy string

const (
	Single     Strategy = "SINGLE"
	Analytical Strategy = "ANALYTICAL "
	Covered    Strategy = "COVERED"
	Vertical   Strategy = "VERTICAL"
	Calender   Strategy = "CALENDAR"
	Strangle   Strategy = "STRANGLE"
	Straddle   Strategy = "STRADDLE"
	Butterfly  Strategy = "BUTTERFLY"
	Condor     Strategy = "CONDOR"
	Diagonal   Strategy = "DIAGONAL"
	Collar     Strategy = "COLLAR"
	Roll       Strategy = "ROLL"
)

type Range string

const (
	ITM Range = "ITM"
	NTM Range = "NTM"
	OTM Range = "OTM"
	SAK Range = "SAK"
	SBK Range = "SBK"
	SNK Range = "SNK"
	ALL Range = "ALL"
)

type ChainsParams struct {
	Symbol           string
	ContractType     string
	StrikeCount      int
	IncludeQuotes    bool
	Strategy         Strategy
	Interval         string
	Strike           string
	Range            Range
	FromDate         Time
	ToDate           Time
	Volatility       float64
	UnderlyingPrice  float64
	InterestRate     float64
	DaysToExpiration int
	ExpMonth         string
	OptionType       string
}

/*func (c ChainsParams) Validate() error {
	for _, status := range c.Status {
		err := status.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}*/

func (c ChainsParams) Params() []string {
	params := []string{fmt.Sprintf("symbol=%s", c.Symbol)}
	params = append(params, fmt.Sprintf("contractType=%s", c.ContractType))
	params = append(params, fmt.Sprintf("strikeCount=%d", c.StrikeCount))
	params = append(params, fmt.Sprintf("includeQuotes=%v", c.IncludeQuotes))
	params = append(params, fmt.Sprintf("strategy=%v", c.Strategy))
	params = append(params, fmt.Sprintf("interval=%s", c.Interval))
	params = append(params, fmt.Sprintf("strike=%s", c.Strike))
	params = append(params, fmt.Sprintf("range=%v", c.Range))

	if !c.FromDate.Time().IsZero() {
		end := c.FromDate.Time().Format(time.RFC3339Nano)
		params = append(params, fmt.Sprintf("fromDate=%s", end))
	}
	if !c.ToDate.Time().IsZero() {
		start := c.ToDate.Time().Format(time.RFC3339Nano)
		params = append(params, fmt.Sprintf("toDate=%s", start))
	}

	return params
}

func (c *TDAmeritrade) GetChains(ctx context.Context, params ChainsParams) (Chains, error) {
	var res Chains
	path := fmt.Sprintf("/%s%s", chains, qsx.Query(params.Params()))
	return res, c.API.Get(ctx, path, &res)
}
