package coinbasepro

import (
	"context"
	"errors"
	"fmt"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
)

// Order
//
// Price Improvement
// Orders are matched against existing order book orders at the price of the order on the book, not at the price of the taker order.
//
// EXAMPLE
// User A places a Buy order for 1 BTC at 100 USD. User B then wishes to sell 1 BTC at 80 USD.
// Because User A's order was first to the trading engine, they will have price priority and the trade will occur at 100 USD.
//
// Order Lifecycle
// Valid orders sent to the matching engine are confirmed immediately and are in the `received` state. If an order executes
// against another order immediately, the order is considered `done`. An order can execute in part or whole. Any part of
// the order not filled immediately will be considered `open`. Orders will stay in the `open` state until canceled or
// subsequently filled by new orders. Orders that are no longer eligible for matching (filled or canceled) are in the done state.
//
// Order Status and Settlement
// Orders which are no longer resting on the order book, will be marked with the `done` status. There is a small window
// between an order being `done` and `settled`. An order is settled when all of the fills have settled and the remaining
// holds (if any) have been removed.
//
// Polling
// For high-volume trading it is strongly recommended that you maintain your own list of open orders and use one of the
// streaming market data feeds to keep it updated. You should poll the open orders endpoint once when you start trading
// to obtain the current state of any open orders.
type Order struct {
	// CreatedAt is the order creation time
	CreatedAt Time `json:"created_at"`
	// CreatedAt is the order completion time
	DoneAt Time `json:"done_at"`
	// DoneReason describes how the order completed
	DoneReason string `json:"done_reason"`
	// ExecutedValue is the value of the order at completion time
	ExecutedValue *float64 `json:"executed_value,string"`
	// FillFees are the total fees for the order
	FillFees *float64 `json:"fill_fees,string"`
	// FilledSize is the amount of Product the filled the order
	FilledSize *float64 `json:"filled_size,string"`
	// Funds is the amount of Funds the the Order actually provides
	Funds *float64 `json:"funds,string"`
	// ID is the server-generated ID of the Order.
	ID string `json:"id"`
	// PostOnly indicates whether only maker orders can be placed. No orders will be matched when post_only mode is active.
	// When PostOnly is true, if any part of the order results in taking liquidity the order will be rejected and no part of it will execute.
	PostOnly bool `json:"post_only"`
	// ProductID identifies the Product associated with the Order
	ProductID ProductID `json:"product_id"`
	// Settled indicates settlement status
	Settled bool `json:"settled"`
	// Side of order, `buy` or `sell`
	Side Side `json:"side"`
	// Size indicates the amount of Product to buy or sell.
	Size float64 `json:"size,string"`
	// SpecifiedFunds TODO: ?
	SpecifiedFunds *float64 `json:"specified_funds,string"`
	// Status indicates the current Order state in the Order Lifecycle
	Status OrderStatus `json:"status"`
	// SelfTradePrevention determines the method(optional).
	// Self-trading is not allowed on Coinbase Pro. Two orders from the same user will not be allowed to match with
	// one another. SelfTradeDecrementAndCancel is the default.
	SelfTradePrevention SelfTrade `json:"stp"`
	// Type of order, `limit` or `market` (`default` is limit)
	Type OrderType `json:"type"`
}

type OrderFilter struct {
	// ProductID limits the list of Orders to those with the specified ProductID
	ProductID ProductID `json:"product-id"`
	// Status limits list of Orders to the provided OrderStatuses. The default, OrderStatusParamAll, returns orders of all statuses
	Status []OrderStatusParam `json:"status"`
}

func (o OrderFilter) Validate() error {
	for _, status := range o.Status {
		err := status.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (o OrderFilter) Params() []string {
	if len(o.Status) == 0 && o.ProductID == "" {
		return nil
	}
	params := make([]string, 0, len(o.Status)+1)
	if o.ProductID != "" {
		params = append(params, fmt.Sprintf("product_id=%s", o.ProductID))
	}
	for _, s := range o.Status {
		params = append(params, fmt.Sprintf("status=%s", s))
	}
	return params
}

// OrderStatus indicates the current Order state in the Order Lifecycle
type OrderStatus string

const (
	OrderStatusActive   OrderStatus = "active"
	OrderStatusDone     OrderStatus = "done"
	OrderStatusOpen     OrderStatus = "open"
	OrderStatusPending  OrderStatus = "pending"
	OrderStatusReceived OrderStatus = "received"
	OrderStatusSettled  OrderStatus = "settled"
)

// OrderStatusParam is a superset of OrderStatus which includes the filter-only OrderStatusParamAll
type OrderStatusParam OrderStatus

const (
	OrderStatusParamActive   = OrderStatusParam(OrderStatusActive)
	OrderStatusParamDone     = OrderStatusParam(OrderStatusDone)
	OrderStatusParamOpen     = OrderStatusParam(OrderStatusOpen)
	OrderStatusParamPending  = OrderStatusParam(OrderStatusPending)
	OrderStatusParamReceived = OrderStatusParam(OrderStatusReceived)
	OrderStatusParamSettled  = OrderStatusParam(OrderStatusSettled)

	// OrderStatusParamAll is only valid as a list filter param
	OrderStatusParamAll OrderStatusParam = "all"
)

func (o OrderStatusParam) Validate() error {
	switch o {
	case OrderStatusParamActive, OrderStatusParamDone, OrderStatusParamOpen, OrderStatusParamPending, OrderStatusParamReceived, OrderStatusParamSettled, OrderStatusParamAll:
		return nil
	default:
		return fmt.Errorf("status(%q) is not valid", o)
	}
}

// LimitOrder is both the default and basic order type. A LimitOrder requires specifying a Price} and Size. The Size
// is the number of base currency to buy or sell, and the Price is the price per base currency. The LimitOrder will be
// filled at the price specified or better. A `sell` order can be filled at the specified price per base currency or a
// higher price per base currency and a `buy` order can be filled at the specified price or a lower price depending on
// market conditions. If market conditions cannot fill the limit order immediately, then the limit order will become
// part of the open order book until filled by another incoming order or canceled by the user.
type LimitOrder struct {
	// Common Order fields
	// ClientOrderID is a client provided order UUID (optional)
	// The received ClientOrderID field value will be broadcast in the public feed for received messages.
	// The ClientOrderID can be used to identify orders in the public feed.
	// If you are consuming the public feed and see a received message with your ClientOrderID,
	// your application *must* record the server-assigned order_idl; the ClientOrderID will not be present on subsequent order status updates.
	// The ClientOrderID will *not* be used after initial broadcast message is sent.
	ClientOrderID string `json:"client_oid"`
	// Must be a valid product ID
	ProductID ProductID `json:"product_id"`
	// SelfTradePrevention determines the method used to prevent orders from the same user from matching with each other.
	// The default SelfTrade prevention method is SelfTradeDecrementAndCancel.
	SelfTradePrevention SelfTrade `json:"stp"`
	// Side type of order, `buy` or `sell`
	Side Side `json:"side"`

	// Stop type of order `entry` or `loss` - requires `stop_price`
	Stop Stop `json:"stop"`
	// StopPrice only if `stop` is defined, sets trigger price for stop order
	StopPrice *float64 `json:"stop_price,string"`
	// Type of order, `limit` or `market` (`default` is limit)

	// LimitOrder must always have an OrderType of `limit`
	Type OrderType `json:"type"`

	// Limit-specific fields
	// TODO: I don't see an example of this, need to make one... WTF does it look like?
	// CancelAfter min, hour, day. (optional) except for LimitOrderTimeInForce_GoodTillTime (required)
	CancelAfter string `json:"cancel_after"`
	// TODO: what does this do? Is it a "dry-run", or some other limitation. Docs suck.
	// PostOnly indicates whether only maker orders can be placed. Invalid when TimeInForce is IOC or FOK.
	// No orders will be matched when PostOnly mode is active.
	// PostOnly indicates that the order should only make liquidity. If any part of the order results in taking liquidity,
	// the order will be rejected and no part of it will execute.
	PostOnly bool `json:"post_only"`
	// Price per unit of product
	// Price must be specified in quote_increment product units. The quote increment is the smallest unit of price.
	// For example, the BTC-USD product has a quote increment of 0.01 or 1 penny.
	// A Price of less than 1 penny will not be accepted, and no fractional penny prices will be accepted.
	Price float64 `json:"price,string"`
	// Size indicates the amount of Product to buy or sell.
	// Size must be greater than the Product BaseMinSize and no larger than the BaseMaxSize.
	// Size can be in incremented in units of BaseIncrement.
	Size float64 `json:"size,string"`
	// TimeInForce describes the time the order is effective: GTC, GTT, IOC, or FOK (default is GTC)
	TimeInForce TimeInForce `json:"time_in_force"`
}

// Validate likely needs more validations
func (l *LimitOrder) Validate() error {
	if l.Type != OrderTypeLimit {
		return errors.New("limit orders must be of type 'limit'")
	}
	if l.ProductID == "" {
		return errors.New("'product_id' is required")
	}
	if l.Price == 0 {
		return errors.New("'price' is required")
	}
	if l.Size == 0 {
		return errors.New("'size' is required")
	}
	if err := l.Side.Validate(); err != nil {
		return err
	}
	if err := l.Stop.Validate(); err != nil {
		return err
	}
	return l.Stop.ValidatePrice(l.StopPrice)
}

// MarketOrder differs from a LimitOrder in that MarketOrder provided no pricing guarantees. MarketOrder do provide a way
// buy or sell specific amounts of base currency or fiat without having to specify the Price. A MarketOrder will execute
// immediately and no part of the order will go on the open order book. Market orders are always considered
// takers and incur taker fees. When placing a market order you can specify Funds and/or Size. Funds will limit how much
// of your quote currency account balance is used and Size will limit the amount of base currency transacted.
// MarketOrder requires one of Funds or Size.
type MarketOrder struct {
	// Common Order fields
	// ClientOrderID is a client provided order UUID (optional)
	// The received ClientOrderID field value will be broadcast in the public feed for received messages.
	// The ClientOrderID can be used to identify orders in the public feed.
	// If you are consuming the public feed and see a received message with your ClientOrderID,
	// your application *must* record the server-assigned order id; the ClientOrderID will not be present on subsequent order status updates.
	// The ClientOrderID will *not* be used after initial broadcast message is sent.
	ClientOrderID string `json:"client_oid"`
	// ProductID must be a valid product ID
	ProductID ProductID `json:"product_id"`
	// SelfTradePrevention determines the method used to prevent orders from the same user from matching with each other.
	// The default SelfTrade prevention method is SelfTradeDecrementAndCancel.
	SelfTradePrevention SelfTrade `json:"stp"`
	// Side type of order, `buy` or `sell`
	Side Side `json:"side"`
	// Stop type of order, `entry` or `loss` - requires `stop_price`
	Stop Stop `json:"stop"`
	// StopPrice only if `stop` is defined, sets trigger price for stop order
	StopPrice *float64 `json:"stop_price,string"`
	// Type of order, `limit` or `market` (`default` is limit)
	// MarketOrder must always have Type `market`
	Type OrderType `json:"type"`

	// MarketOrder-specific fields
	// Funds desired amount of quote currency to use (optional).
	// Funds indicates how much of the product quote currency to buy or sell. For example, a market buy for BTC-USD
	// With Funds of 150.00, the order will use 150 USD to buy BTC (including any fees).
	// Without Funds, the size *must* be specified and Coinbase Pro will use available funds in your account to execute the order.
	// 	A MarkerOrder with Type `sell` can also specify Funds, if provided Funds will limit the `sell` to the amount specified.
	Funds *float64 `json:"funds,string"`
	// Size is number of units of product to buy or sell (optional)
	Size *float64 `json:"size,string"`
}

func (m *MarketOrder) Validate() error {
	if m.Type != OrderTypeMarket {
		return errors.New("market orders must be of type 'market'")
	}
	if m.ProductID == "" {
		return errors.New("'product_id' is required")
	}
	if err := m.Side.Validate(); err != nil {
		return err
	}
	if err := m.Stop.Validate(); err != nil {
		return err
	}
	if err := m.Stop.ValidatePrice(m.StopPrice); err != nil {
		return err
	}
	if m.Funds == nil && m.Size == nil {
		return errors.New("without 'funds', a 'size' is required")
	}
	return nil
}

type SelfTrade string

const (
	// SelfTradeDecrementAndCancel is the default behavior. When two orders from the same user cross, the smaller order
	// will be canceled and the larger order size will be decremented by the smaller order size.
	// If the two orders are the same size, both will be canceled.
	SelfTradeDecrementAndCancel SelfTrade = "dc"
	// SelfTradeCancelNewest cancels the newer (taking) order in full. The old resting order remains on the order book.
	SelfTradeCancelNewest SelfTrade = "cn"
	// SelfTradeCancelOldest cancels the older (resting) order in full. The new order continues to execute.
	SelfTradeCancelOldest SelfTrade = "co"
	// SelfTradeCancelBoth immediately cancels both orders.
	SelfTradeCancelBoth SelfTrade = "cb"

	// Note: When a MarketOrder using the default SelfTradeDecrementAndCancel prevention encounters an open limit order,
	// the behavior depends on the specific values of fields for the MarketOrder message:
	// If Funds and Size are specified for a `buy` order, then the Size of the market order will be decremented internally
	// within the matching engine and Funds will remain unchanged. The intent is to offset target Size without limiting buying power.
	// If Size is not specified, then Funds will be decremented.
	// For a MarketOrder `sell`, the Size will be decremented if there are existing `limit` orders.
)

func (s SelfTrade) Validate() error {
	switch s {
	case SelfTradeDecrementAndCancel, SelfTradeCancelNewest, SelfTradeCancelOldest, SelfTradeCancelBoth:
		return nil
	}
	return fmt.Errorf("stp(%q) is not valid", s)
}

// OrderType influences which other order parameters are required and how the order will be executed by the matching engine.
// If OrderType is not specified, the order will default to a LimitOrder.
type OrderType string

const (
	OrderTypeMarket OrderType = "market"
	OrderTypeLimit  OrderType = "limit"
)

type Side string

const (
	SideBuy  Side = "buy"
	SideSell Side = "sell"
)

func (s Side) Validate() error {
	switch s {
	case SideBuy, SideSell:
		return nil
	default:
		return fmt.Errorf("'side' %q is invalid", s)
	}
}

// Stop orders become active and wait to trigger based on the movement of the last trade price.
// There are two types of triggered Stop, StopLoss and StopEntry:
// StopEntry triggers when the last trade price changes to a value at or above the StopPrice.
// StopLoss triggers when the last trade price changes to a value at or below the StopPrice.
// The last trade price is the last price at which an order was filled. This price can be found in the latest match
// message. Note that not all match messages may be received due to dropped messages.
type Stop string

const (
	// StopEntry triggers when the last trade price changes to a value at or above the StopPrice.
	StopEntry Stop = "entry"
	// StopLoss triggers when the last trade price changes to a value at or below the StopPrice.
	StopLoss Stop = "loss"
	// StopNone is the default Stop. A LimitOrder with StopNone will be booked immediately when created.
	StopNone Stop = ""
)

func (s Stop) Validate() error {
	switch s {
	case StopEntry, StopLoss, StopNone:
		return nil
	default:
		return fmt.Errorf("stop(%q) is not valid", s)
	}
}

func (s Stop) ValidatePrice(price *float64) error {
	switch s {
	case StopEntry, StopLoss:
		if price == nil || *price == 0 || *price < 0 {
			return fmt.Errorf("stop(%q) requires a positive 'stop_price'", s)
		}
	case StopNone:
		if price != nil {
			return errors.New("only stops 'entry' and 'loss' support 'stop_price'")
		}
	default:
		return fmt.Errorf("stop(%q) is not valid", s)
	}
	return nil
}

// LimitOrderSpecific are fields that only apply to an OrderTypeLimit Order, which is both the default and the basic order type.
// A LimitOrder requires specifying a Price and Size. The Size is the number of base currency to buy or sell, and the Price
// is the price per base currency. The LimitOrder will be filled at the price specified or better. A SideSell order can
// be filled at the specified price per base currency or a higher price per base currency and a SideBuy order can be filled
// at the specified price or a lower price depending on market conditions. If market conditions cannot fill the limit order
// immediately, then the limit order will become part of the open order book until filled by another incoming order or
// canceled by the user.
type LimitOrderSpecific struct {
	// TODO: I don't see an example of this, need to make one... WTF does it look like?
	// CancelAfter min, hour, day. (optional) except for TimeInForceGoodTillTime (required)
	CancelAfter string `json:"cancel_after"`
	// TODO: what does this do? Is it a "dry-run", or some other limitation. Docs suck.
	// PostOnly indicates whether only maker orders can be placed. Invalid when time_in_force is IOC or FOK.
	// No orders will be matched when post_only mode is active.
	// PostOnly indicates that the order should only make liquidity. If any part of the order results in taking liquidity,
	// the order will be rejected and no part of it will execute.
	PostOnly bool `json:"post_only"`
	// Price per unit of product
	// Price must be specified in quote_increment product units. The quote increment is the smallest unit of price.
	// For example, the BTC-USD product has a quote increment of 0.01 or 1 penny.
	// A Price of less than 1 penny will not be accepted, and no fractional penny prices will be accepted.
	Price float64 `json:"price"`
	// Size indicates the amount of Product to buy or sell.
	// Size must be greater than the Product BaseMinSize and no larger than the BaseMaxSize.
	// Size can be in incremented in units of BaseIncrement.
	Size float64 `json:"size"`
	// TimeInForce describes the time the order is effective: GTC, GTT, IOC, or FOK (default is GTC)
	TimeInForce TimeInForce `json:"time_in_force"`
}

// TimeInForce policies provide guarantees about the lifetime of an order. There are four policies:
// good till canceled GTC, good till time GTT, immediate or cancel IOC, and fill or kill FOK.
type TimeInForce string

const (
	// TimeInForceGoodTillCanceled orders remain open on the book until canceled. Default.
	TimeInForceGoodTillCanceled TimeInForce = "GTC"
	// TimeInForceGoodTillTime orders remain open on the book until canceled or the allotted CancelAfter has passed.
	// After the CancelAfter timestamp is passed, GTT orders are guaranteed to cancel before any other order is processed.
	// Note that a "day" is considered 24 hours.
	TimeInForceGoodTillTime TimeInForce = "GTT"
	// TimeInForceImmediateOrCancel orders instantly cancel the remaining size of the limit order instead of opening it on the book.
	TimeInForceImmediateOrCancel TimeInForce = "IOC"
	// TimeInForceFillOrKill are rejected if the entire size cannot be matched.
	TimeInForceFillOrKill TimeInForce = "FOK"
)

func (t TimeInForce) Validate() error {
	switch t {
	case TimeInForceGoodTillCanceled, TimeInForceGoodTillTime, TimeInForceImmediateOrCancel, TimeInForceFillOrKill:
		return nil
	default:
		return fmt.Errorf("time_in_force(%q) is not valid", t)
	}
}

func (t TimeInForce) ValidateCancelAfter(cancelAfter string) error {
	if t == TimeInForceGoodTillTime && cancelAfter == "" {
		return fmt.Errorf("time_in_force(%q) requires 'cancel_after'", t)
	}
	return nil
}

// MarketOrderSpecific fields apply only to a MarketOrder. A MarketOrder differs from a LimitOrder in that they provide
// no pricing guarantees. They however do provide a way to buy or sell specific amounts of base currency or fiat without
// having to specify the price. Market orders execute immediately and no part of the market order will go on the open
// order book. Market orders are always considered takers and incur taker fees. When placing a market order you can specify
// funds and/or size. Funds will limit how much of your quote currency account balance is used and Size will limit the
// amount of base currency transacted.
// A MarketOrder requires one of Funds or Size.
type MarketOrderSpecific struct {
	// Funds desired amount of quote currency to use (optional).
	// Funds indicates how much of the product quote currency to buy or sell. For example, a market buy for BTC-USD
	// With Funds of 150.00, the order will use 150 USD to buy BTC (including any fees).
	// Without Funds, the size *must* be specified and Coinbase Pro will use available funds in your account to execute the order.

	// 	A MarkerOrder with Type `sell` can also specify Funds, if provided Funds will limit the `sell` to the amount specified.
	Funds float64 `json:"funds"`
	// Size is number of units of product to buy or sell (optional)
	Size float64 `json:"size"`
}

// Orders is a paged collection of Orders
type Orders struct {
	Orders []*Order    `json:"orders"`
	Page   *Pagination `json:"page,omitempty"`
}

type CancelOrderSpec struct {
	// OrderID or ClientOrderID required.
	OrderID string `json:"order_id"`
	// ClientOrderID or OrderID is required.
	ClientOrderID string `json:"client_oid"`
	// ProductID is optional, but recommended for better performance
	ProductID ProductID `json:"product_id"`
}

func (c CancelOrderSpec) Validate() error {
	if (c.OrderID != "" && c.ClientOrderID != "") || (c.OrderID == "" && c.ClientOrderID == "") {
		return errors.New("one and only one of 'order_id' or 'client_oid' must be provided")
	}
	return nil
}

func (c CancelOrderSpec) Params() []string {
	if c.ProductID != "" {
		return []string{fmt.Sprintf("product_id=%s", c.ProductID)}
	}
	return nil
}

// CreateLimitOrder creates a LimitOrder to trade a Product with specified Price and Size limits.
func (c *CoinbasePro) CreateLimitOrder(ctx context.Context, limitOrder LimitOrder) (Order, error) {
	if err := limitOrder.Validate(); err != nil {
		return Order{}, err
	}
	var order Order
	path := fmt.Sprintf("/%s/", coinbaseproOrders)
	return order, c.API.Post(ctx, path, limitOrder, &order)
}

// CreateMarketOrder creates a MarketOrder with no pricing guarantees. A MarketOrder makes it easy to trade specific
// amounts of a Product without specifying prices.
func (c *CoinbasePro) CreateMarketOrder(ctx context.Context, marketOrder MarketOrder) (Order, error) {
	if err := marketOrder.Validate(); err != nil {
		return Order{}, err
	}
	var order Order
	path := fmt.Sprintf("/%s/", coinbaseproOrders)
	return order, c.API.Post(ctx, path, marketOrder, &order)
}

// CancelOrder cancels a previously placed order. orderID is mandatory, productID is optional but will make the request
// more performant. If the Order had no matches during its lifetime, it may be subject to purge and as a result will
// no longer available via GetOrder.
// Requires "trade" permission.
func (c *CoinbasePro) CancelOrder(ctx context.Context, spec CancelOrderSpec) (map[string]interface{}, error) {
	if err := spec.Validate(); err != nil {
		return nil, err
	}
	var resp map[string]interface{}
	err := spec.Validate()
	if err != nil {
		return nil, err
	}
	//path := fmt.Sprintf("/%s/", coinbaseproOrders)
	return resp, c.API.Delete(ctx, "/orders/"+qsx.Query(spec.Params()), nil, &resp)
}

// GetOrders retrieves a paginated list of the current open orders for the current Profile. Only open or un-settled
// orders are returned by default. An OrderFilter can be used to further refine the request.
func (c *CoinbasePro) GetOrders(ctx context.Context, filter OrderFilter, pagination PaginationParams) (Orders, error) {
	params := append(filter.Params(), pagination.Params()...)
	var orders Orders
	path := fmt.Sprintf("/%s/%s", coinbaseproOrders, qsx.Query(params))
	return orders, c.API.Get(ctx, path, &orders)
}

// GetOrder retrieves the details of a single Order. The requested Order must belong to the current Profile.
func (c *CoinbasePro) GetOrder(ctx context.Context, orderID string) (Order, error) {
	var order Order
	path := fmt.Sprintf("/%s/%s", coinbaseproOrders, orderID)
	return order, c.API.Get(ctx, path, &order)
}

// GetClientOrder retrieves the details of a single Order using a client-provided identifier.
// The requested Order must belong to the current Profile.
func (c *CoinbasePro) GetClientOrder(ctx context.Context, clientID string) (Order, error) {
	var order Order
	//path := fmt.Sprintf("/%s/%s", coinbaseproOrders, orderID)
	return order, c.API.Get(ctx, fmt.Sprintf("/orders/client:%s", clientID), &order)
}
