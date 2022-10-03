package qsx

// Product
// Only a maximum of one of trading_disabled, cancel_only, post_only, limit_only can be true at once. If none are true,
// the product is trading normally.
// !! When limit_only is true, matching can occur if a limit order crosses the book.
// !! Product ID will not change once assigned to a Product but all other fields are subject to change.
type Product struct {
	ID string `json:"id"`

	// BaseCurrency is the base in the pair of currencies that comprise the Product
	BaseCurrency string `json:"base_currency"`

	// QuoteCurrency
	QuoteCurrency string `json:"quote_currency"`

	// BaseMinSize defines the minimum order size
	BaseMinSize string `json:"base_min_size"`

	// BaseMaxSize defines the maximum order size
	BaseMaxSize string `json:"base_max_size"`

	// QuoteIncrement
	QuoteIncrement string `json:"quote_increment"`

	// BaseIncrement specifies the minimum increment for the BaseCurrency
	BaseIncrement string `json:"base_increment"`

	// DisplayName
	DisplayName string `json:"display_name"`

	// MinMarketFunds defines the minimum funds allowed
	MinMarketFunds string `json:"min_market_funds"`

	// MaxMarketFunds defines the maximum funds allowed
	MaxMarketFunds string `json:"max_market_funds"`

	// MarginEnabled
	MarginEnabled bool `json:"margin_enabled"`

	// PostOnly indicates whether only maker orders can be placed. No orders will be matched when post_only mode is active.
	PostOnly bool `json:"post_only"`

	// LimitOnly indicates whether this product only accepts limit orders.
	LimitOnly bool `json:"limit_only"`

	// CancelOnly indicates whether this product only accepts cancel requests for orders.
	CancelOnly bool `json:"cancel_only"`

	// Status
	Status string `json:"status"`

	// StatusMessage provides any extra information regarding the status, if available
	StatusMessage string `json:"status_message"`

	// AuctionMode
	AuctionMode bool `json:"auction_mode"`
}
