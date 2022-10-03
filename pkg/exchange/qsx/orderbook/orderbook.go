package orderbook

import (
	"fmt"
	"log"
	"sync"
)

// MaxLimitsNum maximum limits per order book side to pre-allocate memory
const MaxLimitsNum int = 10000

type Orderbook struct {
	Bids           *Tree
	Asks           *Tree
	bidLimitsCache map[float64]*LimitOrder
	askLimitsCache map[float64]*LimitOrder
	pool           *sync.Pool
}

func NewOrderbook() *Orderbook {
	bids := NewTree()
	asks := NewTree()
	return &Orderbook{
		Bids:           &bids,
		Asks:           &asks,
		bidLimitsCache: make(map[float64]*LimitOrder, MaxLimitsNum),
		askLimitsCache: make(map[float64]*LimitOrder, MaxLimitsNum),
		pool: &sync.Pool{
			New: func() interface{} {
				limit := NewLimitOrder(0.0)
				return &limit
			},
		},
	}
}

func (book *Orderbook) Add(price float64, o *Order) {
	var limit *LimitOrder

	if o.BidOrAsk {
		limit = book.bidLimitsCache[price]
	} else {
		limit = book.askLimitsCache[price]
	}

	if limit == nil {
		// getting a new limit from pool
		limit = book.pool.Get().(*LimitOrder)
		limit.Price = price

		// insert into the corresponding BST and cache
		if o.BidOrAsk {
			book.Bids.Put(price, limit)
			book.bidLimitsCache[price] = limit
		} else {
			book.Asks.Put(price, limit)
			book.askLimitsCache[price] = limit
		}
	}

	// add order to the limit
	limit.Enqueue(o)
}

func (book *Orderbook) NewOrder() *Order {
	return &Order{}
}

func (book *Orderbook) Cancel(o *Order) {
	limit := o.Limit
	limit.Delete(o)

	if limit.Size() == 0 {
		// remove the limit if there are no orders
		if o.BidOrAsk {
			book.Bids.Delete(limit.Price)
			delete(book.bidLimitsCache, limit.Price)
		} else {
			book.Asks.Delete(limit.Price)
			delete(book.askLimitsCache, limit.Price)
		}

		// put it back to the pool
		book.pool.Put(limit)
	}
}

func (book *Orderbook) ClearBidLimit(price float64) {
	book.clearLimit(price, true)
}

func (book *Orderbook) ClearAskLimit(price float64) {
	book.clearLimit(price, false)
}

func (book *Orderbook) clearLimit(price float64, bidOrAsk bool) {
	var limit *LimitOrder
	if bidOrAsk {
		limit = book.bidLimitsCache[price]
	} else {
		limit = book.askLimitsCache[price]
	}

	if limit == nil {
		log.Panicln(fmt.Sprintf("there is no such price limit %0.8f", price))
	}

	limit.Clear()
}

func (book *Orderbook) DeleteBidLimit(price float64) {
	limit := book.bidLimitsCache[price]
	if limit == nil {
		return
	}

	book.deleteLimit(price, true)
	delete(book.bidLimitsCache, price)

	// put limit back to the pool
	limit.Clear()
	book.pool.Put(limit)

}

func (book *Orderbook) DeleteAskLimit(price float64) {
	limit := book.bidLimitsCache[price]
	if limit == nil {
		return
	}

	book.deleteLimit(price, false)
	delete(book.askLimitsCache, price)

	// put limit back to the pool
	limit.Clear()
	book.pool.Put(limit)
}

func (book *Orderbook) deleteLimit(price float64, bidOrAsk bool) {
	if bidOrAsk {
		book.Bids.Delete(price)
	} else {
		book.Asks.Delete(price)
	}
}

func (book *Orderbook) GetVolumeAtBidLimit(price float64) float64 {
	limit := book.bidLimitsCache[price]
	if limit == nil {
		return 0
	}
	return limit.TotalVolume()
}

func (book *Orderbook) GetVolumeAtAskLimit(price float64) float64 {
	limit := book.askLimitsCache[price]
	if limit == nil {
		return 0
	}
	return limit.TotalVolume()
}

func (book *Orderbook) GetBestBid() float64 {
	return book.Bids.Max()
}

func (book *Orderbook) GetBestOffer() float64 {
	return book.Asks.Min()
}

func (book *Orderbook) BLength() int {
	return len(book.bidLimitsCache)
}

func (book *Orderbook) ALength() int {
	return len(book.askLimitsCache)
}
