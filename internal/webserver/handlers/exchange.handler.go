package handlers

/*
import (
	"encoding/json"
	"github.com/quantstop/quantstopterminal/internal"
	"github.com/quantstop/quantstopterminal/internal/database/models"
	"github.com/quantstop/quantstopterminal/internal/webserver/errors"
	"github.com/quantstop/quantstopterminal/internal/webserver/write"
	"net/http"
)

type setExchangeRequest struct {
	Name           string `json:"name"`
	AuthKey        string `json:"authKey"`
	AuthPassphrase string `json:"authPassphrase"`
	AuthSecret     string `json:"AuthSecret"`
	Currency       string `json:"currency"`
}

func SetExchange(bot internal.IEngine, user *models.User, w http.ResponseWriter, r *http.Request) http.HandlerFunc {

	db, _ := bot.GetSQL("core")

	decoder := json.NewDecoder(r.Body)
	req := setExchangeRequest{}
	err := decoder.Decode(&req)
	if err != nil || &req == nil {
		return write.Error(errors.NoJSONBody)
	}

	if req.Name == "" || req.Currency == "" {
		return write.Error(errors.InvalidInput)
	}

	exchange := &models.Exchange{
		Name:           req.Name,
		UserDefined:    1,
		AuthKey:        req.AuthKey,
		AuthPassphrase: req.AuthPassphrase,
		AuthSecret:     req.AuthSecret,
		Currency:       req.Currency,
	}

	//todo: encrypt api keys ...

	err = exchange.CreateExchange(db)
	if err != nil {
		//todo: can we get more specific errors? do we even need to?
		return write.Error(err)
	}

	return write.Success()
}

type getExchangesResponse struct {
	Type      string              `json:"type"`
	Exchanges []SupportedExchange `json:"exchanges"`
}

type SupportedExchange struct {
	ID string `json:"id"`
}

func GetExchanges(bot internal.IEngine, user *models.User, w http.ResponseWriter, r *http.Request) http.HandlerFunc {

	var supportedExchanges []SupportedExchange
	for _, e := range bot.GetSupportedExchangesList() {
		supportedExchanges = append(supportedExchanges, SupportedExchange{ID: e})
	}

	res := getExchangesResponse{
		Type:      "getExchanges",
		Exchanges: supportedExchanges,
	}

	return write.JSON(res)
}

*/
