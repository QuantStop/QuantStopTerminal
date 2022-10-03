package orderbook

// Order defines a single order in the orderbook, as a node in a LimitOrder FIFO queue
type Order struct {
	Id       int
	Volume   float64
	Next     *Order
	Prev     *Order
	Limit    *LimitOrder
	BidOrAsk bool
}
