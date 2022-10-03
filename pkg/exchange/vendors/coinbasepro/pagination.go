package coinbasepro

import (
	"errors"
	"fmt"
)

// Pagination
// Docs state that "Coinbase Pro uses cursor pagination for all REST requests which return arrays."
// This is *not* true: there are several (/accounts/, /coinbase-accounts/, /currencies/, ...) that return arrays and no pagination.
// Docs (paraphrased) say: Some endpoints (/trades/, /fills/, /orders/, ...) return the latest items by default, along
// with pagination tokens in the header. To retrieve additional results, subsequent requests should provide a
// direction/token using a token copied from the headers of a previous response.
//
// In this implementation, endpoints that support Pagination will return a pluralized struct (Orders, Holds, Ledger, ...)
// and methods that return paginated structs will have a Get prefix (GetOrders, GetHolds, GetLedger, ...). Methods that
// return non-paginated slices of structs will have a List prefix (ListAccounts, ListCoinbaseAccounts, ListCurrencies).
//
// Be warned that Pagination in the coinbase pro REST API is weird. Other client implementations have simplified the bidirectional
// navigation into a unidirectional Cursor implementation. I chose not to do this, which might prove to be a mistake.
//
// Also note that the actual value of the pagination parameter should be considered meaningless/opaque,
// the coinbasepro REST API is free to put whatever it wants in those tokens.

// PaginationParams allow for requests to limit the number of responses as well as request the next page
// of responses either Before or After the current Page in sort order. The sort order is endpoint dependent, though will
// tend to be chronological, with newest first.
// A Limit of 0 will be defaulted to 100, and 100 is also the maximum allowed value.
type PaginationParams struct {
	After  string `json:"after"`
	Before string `json:"before"`
	Limit  int    `json:"limit"`
}

func (p PaginationParams) Params() []string {
	var params []string
	if p.After != "" {
		params = append(params, fmt.Sprintf("after=%s", p.After))
	}
	if p.Before != "" {
		params = append(params, fmt.Sprintf("before=%s", p.Before))
	}
	if p.Limit != 0 {
		params = append(params, fmt.Sprintf("limit=%d", p.Limit))
	}
	return params
}

func (p PaginationParams) Validate() error {
	if p.Before != "" && p.After != "" {
		return errors.New("only one of 'before' or 'after' allowed")
	}
	if p.Limit < 0 || p.Limit > 100 {
		return errors.New("limit %d is outside of allowed range (0..100)")
	}
	return nil
}

// Pagination uses the response from the coinbase pro REST API to populate a pair of page tokens which can be used to
// navigate forwards/backwards through a paginated response. The before/after nomenclature is a direct lift from the coinbase pro REST API
// which, as their own docs state, "is unintuitive at first," because most paginated requests return the latest information
// (newest) as the first page. Because the pages are sorted by newest (in chronological time) first, to get older information,
// use the `before` token.
// Just need to keep in mind: the `before` and `after` refer to sort order and not chronological time. It makes sense, but just barely.
type Pagination struct {
	After  string `json:"after,omitempty"`
	Before string `json:"before,omitempty"`
}

// NotEmpty provides a convenience for determining when a response is Paginated. Pagination will be populated for Paginated
// responses.
func (p *Pagination) NotEmpty() bool {
	return !(p.After == "" && p.Before == "")
}
