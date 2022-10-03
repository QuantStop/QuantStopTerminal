package orderbook

import "log"

// LimitOrder price orders combined as a FIFO queue
type LimitOrder struct {
	Price       float64
	orders      *OrdersQueue
	totalVolume float64
}

func NewLimitOrder(price float64) LimitOrder {
	q := NewOrdersQueue()
	return LimitOrder{
		Price:  price,
		orders: &q,
	}
}

func (lo *LimitOrder) TotalVolume() float64 {
	return lo.totalVolume
}

func (lo *LimitOrder) Size() int {
	return lo.orders.Size()
}

func (lo *LimitOrder) Enqueue(o *Order) {
	lo.orders.Enqueue(o)
	o.Limit = lo
	lo.totalVolume += o.Volume
}

func (lo *LimitOrder) Dequeue() *Order {
	if lo.orders.IsEmpty() {
		return nil
	}

	o := lo.orders.Dequeue()
	lo.totalVolume -= o.Volume
	return o
}

func (lo *LimitOrder) Delete(o *Order) {
	if o.Limit != lo {
		log.Panicln("order does not belong to the limit")
	}

	lo.orders.Delete(o)
	o.Limit = nil
	lo.totalVolume -= o.Volume
}

func (lo *LimitOrder) Clear() {
	q := NewOrdersQueue()
	lo.orders = &q
	lo.totalVolume = 0
}
