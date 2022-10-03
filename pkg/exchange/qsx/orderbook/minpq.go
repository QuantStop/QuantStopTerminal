package orderbook

import "log"

// MinPQ Minimum Oriented Priority Queue
type MinPQ struct {
	keys []float64
	n    int
}

func NewMinPQ(size int) MinPQ {
	return MinPQ{
		keys: make([]float64, size+1),
	}
}

func (pq *MinPQ) Size() int {
	return pq.n
}

func (pq *MinPQ) IsEmpty() bool {
	return pq.n == 0
}

func (pq *MinPQ) Insert(key float64) {
	if pq.n+1 == cap(pq.keys) {
		log.Panicln("pq is full")
	}

	pq.n++
	pq.keys[pq.n] = key

	// restore order: LogN
	pq.swim(pq.n)
}

func (pq *MinPQ) Top() float64 {
	if pq.IsEmpty() {
		log.Panicln("pq is empty")
	}

	return pq.keys[1]
}

// DelTop removes minimal element and returns it
func (pq *MinPQ) DelTop() float64 {
	if pq.IsEmpty() {
		log.Panicln("pq is empty")
	}

	top := pq.keys[1]
	pq.keys[1] = pq.keys[pq.n]
	pq.n--

	// restore order: logN
	pq.sink(1)

	return top
}

func (pq *MinPQ) swim(k int) {
	for k > 1 && pq.keys[k] < pq.keys[k/2] {
		// swap
		pq.keys[k], pq.keys[k/2] = pq.keys[k/2], pq.keys[k]
		k = k / 2
	}
}

func (pq *MinPQ) sink(k int) {
	for 2*k <= pq.n {
		c := 2 * k
		// select minimum of two children
		if c < pq.n && pq.keys[c+1] < pq.keys[c] {
			c++
		}

		if pq.keys[c] < pq.keys[k] {
			// swap
			pq.keys[c], pq.keys[k] = pq.keys[k], pq.keys[c]
			k = c
		} else {
			break
		}
	}
}
