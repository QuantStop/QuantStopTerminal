package orderbook

// OrdersQueue Doubly linked orders queue
// TODO: this should be compared with ring buffer queue performance
type OrdersQueue struct {
	head *Order
	tail *Order
	size int
}

func NewOrdersQueue() OrdersQueue {
	return OrdersQueue{}
}

func (q *OrdersQueue) Size() int {
	return q.size
}

func (q *OrdersQueue) IsEmpty() bool {
	return q.size == 0
}

func (q *OrdersQueue) Enqueue(o *Order) {
	tail := q.tail
	q.tail = o
	if tail != nil {
		tail.Next = o
	}
	if q.head == nil {
		q.head = o
	}
	q.size++
}

func (q *OrdersQueue) Dequeue() *Order {
	if q.size == 0 {
		return nil
	}

	head := q.head
	if q.tail == q.head {
		q.tail = nil
	}

	q.head = q.head.Next
	q.size--
	return head
}

func (q *OrdersQueue) Delete(o *Order) {
	prev := o.Prev
	next := o.Next
	if prev != nil {
		prev.Next = next
	}
	if next != nil {
		next.Prev = prev
	}
	o.Next = nil
	o.Prev = nil

	q.size--

	if q.head == o {
		q.head = next
	}
	if q.tail == o {
		q.tail = prev
	}
}
