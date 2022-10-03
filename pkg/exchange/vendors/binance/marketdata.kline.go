package binance

import (
	"context"
	"fmt"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
	"time"
)

type KLineResponse []Kline

type Kline struct {
	OpenTime                 time.Time `json:"openTime"`
	Open                     string    `json:"open"`
	High                     string    `json:"high"`
	Low                      string    `json:"low"`
	Close                    string    `json:"close"`
	Volume                   string    `json:"volume"`
	CloseTime                time.Time `json:"closeTime"`
	QuoteAssetVolume         string    `json:"quoteAssetVolume"`
	TradeNum                 int64     `json:"tradeNum"`
	TakerBuyBaseAssetVolume  string    `json:"takerBuyBaseAssetVolume"`
	TakerBuyQuoteAssetVolume string    `json:"takerBuyQuoteAssetVolume"`
}

type KLineParameters struct {
	Symbol    string    `json:"symbol"`
	Interval  string    `json:"interval"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Limit     int       `json:"limit"`
}

func (h *KLineParameters) Params() []string {
	params := []string{fmt.Sprintf("symbol=%s", h.Symbol)}
	params = append(params, fmt.Sprintf("interval=%s", h.Interval))
	params = append(params, fmt.Sprintf("limit=%v", h.Limit))
	if !h.EndTime.IsZero() {
		end := h.EndTime.Format(time.RFC3339Nano)
		params = append(params, fmt.Sprintf("endTime=%s", end))
	}
	if !h.StartTime.IsZero() {
		start := h.StartTime.Format(time.RFC3339Nano)
		params = append(params, fmt.Sprintf("startTime=%s", start))
	}
	return params
}

func (b *Binance) GetKlineData(ctx context.Context, params KLineParameters) (KLineResponse, error) {
	var res KLineResponse
	path := fmt.Sprintf("/%s/%s", candleStick, qsx.Query(params.Params()))
	if err := b.API.Get(ctx, path, &res); err != nil {
		return res, err
	}
	return res, nil
}
