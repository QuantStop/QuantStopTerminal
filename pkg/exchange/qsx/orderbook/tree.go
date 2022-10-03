package orderbook

import (
	"fmt"
	"log"
)

// A self-balancing Binary Search Tree with 2*logN worst case guarantees for
// search, put, delete, min, max, select, rank, floor, ceiling operations.
// Average runtime for search-based operations estimated as 1*logN

type Node struct {
	Key   float64
	Value *LimitOrder
	Next  *Node
	Prev  *Node
	left  *Node
	right *Node
	size  int
	isRed bool
}

type Tree struct {
	root *Node
	minC *Node // cached min/max keys for O(1) access
	maxC *Node
}

func NewTree() Tree {
	return Tree{}
}

func (t *Tree) Size() int {
	return t.size(t.root)
}

func (t *Tree) size(n *Node) int {
	if n == nil {
		return 0
	}
	return n.size
}

func (t *Tree) IsEmpty() bool {
	return t.size(t.root) == 0
}

func (t *Tree) panicIfEmpty() {
	if t.IsEmpty() {
		log.Panicln("Red Black BST is empty")
	}
}

func (t *Tree) Contains(key float64) bool {
	return t.get(t.root, key) != nil
}

func (t *Tree) Get(key float64) *LimitOrder {
	t.panicIfEmpty()

	x := t.get(t.root, key)
	if x == nil {
		log.Panicln(fmt.Sprintf("key %0.8f does not exist", key))
	}

	return x.Value
}

func (t *Tree) get(n *Node, key float64) *Node {
	if n == nil {
		return nil
	}

	if n.Key == key {
		return n
	}

	if n.Key > key {
		return t.get(n.left, key)
	} else {
		return t.get(n.right, key)
	}
}

func (t *Tree) isRed(n *Node) bool {
	if n == nil {
		// nil nodes are black by default
		return false
	}

	return n.isRed
}

func (t *Tree) flipColors(n *Node) {
	if n == nil {
		return
	}

	// inverse children colors
	if n.left != nil {
		n.left.isRed = !n.left.isRed
	}
	if n.right != nil {
		n.right.isRed = !n.right.isRed
	}

	// inverse node color
	n.isRed = !n.isRed
}

func (t *Tree) rotateLeft(n *Node) *Node {
	x := n.right
	n.right = x.left
	x.left = n

	x.isRed = n.isRed
	n.isRed = true

	// re-calculate sizes
	n.size = t.size(n.left) + 1 + t.size(n.right)
	x.size = t.size(x.left) + 1 + t.size(x.right)

	return x
}

func (t *Tree) rotateRight(n *Node) *Node {
	x := n.left
	n.left = x.right
	x.right = n

	x.isRed = n.isRed
	n.isRed = true

	// re-calculate sizes
	n.size = t.size(n.left) + 1 + t.size(n.right)
	x.size = t.size(x.left) + 1 + t.size(x.right)

	return x
}

func (t *Tree) Put(key float64, value *LimitOrder) {
	t.root = t.put(t.root, key, value)

	// keeping root black
	t.root.isRed = false
}

func (t *Tree) put(newNode *Node, key float64, value *LimitOrder) *Node {
	if newNode == nil {
		// search miss, creating a new node with a red link as a part of 3- or 4-node
		newNode = &Node{
			Value: value,
			Key:   key,
			size:  1,
			isRed: true,
		}

		if t.minC == nil || key < t.minC.Key {
			// new min
			t.minC = newNode
		}
		if t.maxC == nil || key > t.maxC.Key {
			// new max
			t.maxC = newNode
		}

		return newNode
	}

	if newNode.Key == key {
		// search hit, updating the value
		newNode.Value = value
		return newNode
	}

	if newNode.Key > key {
		left := newNode.left
		newNode.left = t.put(newNode.left, key, value)
		if left == nil {
			// new node has been just inserted to the left
			prev := newNode.Prev
			if prev != nil {
				prev.Next = newNode.left
			}
			newNode.left.Prev = prev
			newNode.left.Next = newNode
			newNode.Prev = newNode.left
		}
	} else {
		right := newNode.right
		newNode.right = t.put(newNode.right, key, value)
		if right == nil {
			// new node has been just inserted to the right
			next := newNode.Next
			if next != nil {
				next.Prev = newNode.right
			}
			newNode.right.Next = next
			newNode.right.Prev = newNode
			newNode.Next = newNode.right
		}
	}

	// balancing the tree
	if t.isRed(newNode.right) && !t.isRed(newNode.left) {
		// fixing right leaning red link case
		// this can lead to the next case in upper level
		newNode = t.rotateLeft(newNode)
	}
	if t.isRed(newNode.left) && t.isRed(newNode.left.left) {
		// making 4-node
		newNode = t.rotateRight(newNode)
	}
	if t.isRed(newNode.left) && t.isRed(newNode.right) {
		// convert 4-node into 3 2-nodes
		t.flipColors(newNode)
	}

	// re-calc size
	newNode.size = t.size(newNode.left) + 1 + t.size(newNode.right)
	return newNode
}

func (t *Tree) Height() int {
	if t.IsEmpty() {
		return 0
	}

	return t.height(t.root)
}

func (t *Tree) height(n *Node) int {
	if n == nil {
		return 0
	}

	lheight := t.height(n.left)
	rheight := t.height(n.right)

	height := lheight
	if rheight > lheight {
		height = rheight
	}

	return height + 1
}

func (t *Tree) IsRedBlack() bool {
	balanced, _ := t.isBalanced(t.root)
	return balanced && t.is23(t.root)
}

func (t *Tree) isBalanced(n *Node) (bool, int) {
	if n == nil {
		// nil node is black by default
		return true, 1
	}

	lb, l := t.isBalanced(n.left)
	rb, r := t.isBalanced(n.right)

	b := l
	if r > l {
		b = r
	}

	if !t.isRed(n) {
		b += 1
	}

	return lb && rb && l == r, b
}

func (t *Tree) is23(n *Node) bool {
	if n == nil {
		return true
	}

	if t.isRed(n.right) {
		// it should have only left leaning red links
		return false
	}

	if t.isRed(n) && t.isRed(n.left) {
		// no node should be connected by two red links
		return false
	}

	return t.is23(n.left) && t.is23(n.right)
}

func (t *Tree) Min() float64 {
	t.panicIfEmpty()
	return t.minC.Key
}

func (t *Tree) MinValue() *LimitOrder {
	t.panicIfEmpty()
	return t.minC.Value
}

func (t *Tree) MinPointer() *Node {
	t.panicIfEmpty()
	return t.minC
}

func (t *Tree) min(n *Node) *Node {
	if n.left == nil {
		return n
	}

	return t.min(n.left)
}

func (t *Tree) Max() float64 {
	t.panicIfEmpty()
	return t.maxC.Key
}

func (t *Tree) MaxValue() *LimitOrder {
	t.panicIfEmpty()
	return t.maxC.Value
}

func (t *Tree) MaxPointer() *Node {
	t.panicIfEmpty()
	return t.maxC
}

func (t *Tree) max(n *Node) *Node {
	if n.right == nil {
		return n
	}

	return t.max(n.right)
}

func (t *Tree) Floor(key float64) float64 {
	t.panicIfEmpty()

	floor := t.floor(t.root, key)
	if floor == nil {
		log.Panicln(fmt.Sprintf("there are no keys <= %0.8f", key))
	}

	return floor.Key
}

func (t *Tree) floor(n *Node, key float64) *Node {
	if n == nil {
		// search miss
		return nil
	}

	if n.Key == key {
		// search hit
		return n
	}

	if n.Key > key {
		// floor must be in the left sub-tree
		return t.floor(n.left, key)
	}

	// key could be in the right sub-tree, if not, using current root
	floor := t.floor(n.right, key)
	if floor != nil {
		return floor
	}

	return n
}

func (t *Tree) Ceiling(key float64) float64 {
	t.panicIfEmpty()

	ceiling := t.ceiling(t.root, key)
	if ceiling == nil {
		log.Panicln(fmt.Sprintf("there are no keys >= %0.8f", key))
	}

	return ceiling.Key
}

func (t *Tree) ceiling(n *Node, key float64) *Node {
	if n == nil {
		// search miss
		return nil
	}

	if n.Key == key {
		// search hit
		return n
	}

	if n.Key < key {
		// ceiling must be in the right sub-tree
		return t.ceiling(n.right, key)
	}

	// the key could be in the left sub-tree, if not, using current root
	ceiling := t.ceiling(n.left, key)
	if ceiling != nil {
		return ceiling
	}

	return n
}

func (t *Tree) Select(k int) float64 {
	if k < 0 || k >= t.Size() {
		log.Panicln("index out of range")
	}

	return t.selectNode(t.root, k).Key
}

func (t *Tree) selectNode(n *Node, k int) *Node {
	if t.size(n.left) == k {
		return n
	}

	if t.size(n.left) > k {
		return t.selectNode(n.left, k)
	}

	k = k - t.size(n.left) - 1
	return t.selectNode(n.right, k)
}

func (t *Tree) Rank(key float64) int {
	t.panicIfEmpty()
	return t.rank(t.root, key)
}

func (t *Tree) rank(n *Node, key float64) int {
	if n == nil {
		return 0
	}

	if n.Key == key {
		return t.size(n.left)
	}

	if n.Key > key {
		return t.rank(n.left, key)
	}

	return t.size(n.left) + 1 + t.rank(n.right, key)
}

func (t *Tree) moveRedLeft(n *Node) *Node {
	// assuming that n.left and n.left.left are black and n is red,
	// make h.left or one of its children red
	t.flipColors(n)
	// now n is black and both left and right are red

	// fixing red black invariat that no node can be connected with two red links
	if t.isRed(n.right.left) {
		n.right = t.rotateRight(n.right)
		// now n.right and n.right.right are red, fixing that by rotating n
		n = t.rotateLeft(n)
		// now n.right, n.right.right and n.left are red

		t.flipColors(n)
		// now n.left and n.right are black, n.left.left is red
	}

	return n
}

func (t *Tree) DeleteMin() {
	t.panicIfEmpty()

	if !t.isRed(t.root.left) && !t.isRed(t.root.right) {
		// making root red temporarily to fit invariant required for moveRedLeft method
		t.root.isRed = true
	}
	t.root = t.deleteMin(t.root)
	if !t.IsEmpty() {
		t.root.isRed = false
	}
}

func (t *Tree) deleteMin(n *Node) *Node {
	if n.left == nil {
		// we've reached the least leave of the tree
		next := n.Next
		prev := n.Prev
		if prev != nil {
			prev.Next = next
		}
		if next != nil {
			next.Prev = prev
		}
		n.Next = nil
		n.Prev = nil

		// updating global min
		if t.minC == n {
			t.minC = next
		}

		return n.right
	}

	// making current node a part of 3 or 4 node by moving red link to the left
	if !t.isRed(n.left) && !t.isRed(n.left.left) {
		n = t.moveRedLeft(n)
	}

	n.left = t.deleteMin(n.left)

	// we have to restore balance of the tree moving from bottom to top now
	if t.isRed(n.right) {
		n = t.rotateLeft(n)
	}
	if t.isRed(n.left) && t.isRed(n.left.left) {
		n = t.rotateRight(n)
	}
	if t.isRed(n.left) && t.isRed(n.right) {
		t.flipColors(n)
	}

	n.size = t.size(n.left) + 1 + t.size(n.right)
	return n
}

func (t *Tree) moveRedRight(n *Node) *Node {
	// assuming n is red, n.right and n.right.left are black,
	// make h.right or one of its children red
	t.flipColors(n)
	// now n is black and n.right is red

	if t.isRed(n.left.left) {
		// meaning n.left should be red now after fliping the colors of n
		n = t.rotateRight(n)
		// now n.left is red, n.right and n.right.right are red

		t.flipColors(n)
		// now n.left is black, n.right is black, n.right.right is red
	}
	return n
}

func (t *Tree) DeleteMax() {
	t.panicIfEmpty()

	if !t.isRed(t.root.left) && !t.isRed(t.root.right) {
		t.root.isRed = true
	}
	t.root = t.deleteMax(t.root)
	if !t.IsEmpty() {
		t.root.isRed = false
	}
}

func (t *Tree) deleteMax(n *Node) *Node {
	if t.isRed(n.left) {
		// making right red by rotating
		n = t.rotateRight(n)
	}
	if n.right == nil {
		// we've reached the largest key in the tree
		next := n.Next
		prev := n.Prev
		if prev != nil {
			prev.Next = next
		}
		if next != nil {
			next.Prev = prev
		}
		n.Next = nil
		n.Prev = nil

		// updating global max
		if t.maxC == n {
			t.maxC = prev
		}

		return n.left
	}

	// making right left on the way from top to bottom
	if !t.isRed(n.right) && !t.isRed(n.right.left) {
		n = t.moveRedRight(n)
	}

	n.right = t.deleteMax(n.right)

	// balancing back on the way from bottom to top
	if t.isRed(n.right) {
		n = t.rotateLeft(n)
	}
	if t.isRed(n.left) && t.isRed(n.left.left) {
		n = t.rotateRight(n)
	}
	if t.isRed(n.left) && t.isRed(n.right) {
		t.flipColors(n)
	}

	n.size = t.size(n.left) + 1 + t.size(n.right)
	return n
}

func (t *Tree) Delete(key float64) {
	t.panicIfEmpty()

	if !t.isRed(t.root.left) && !t.isRed(t.root.right) {
		t.root.isRed = true
	}
	t.root = t.delete(t.root, key)
	if !t.IsEmpty() {
		t.root.isRed = false
	}
}

func (t *Tree) delete(n *Node, key float64) *Node {
	if n.Key > key {
		if n.left == nil {
			// search miss
			return nil
		}

		// looking into the left sub-tree
		if !t.isRed(n.left) && !t.isRed(n.left.left) {
			n = t.moveRedLeft(n)
		}

		n.left = t.delete(n.left, key)

	} else {
		// checking current node and right sub-tree if required
		if t.isRed(n.left) {
			n = t.rotateRight(n)
		}
		if n.Key == key && n.right == nil {
			// search hit and we don't have right sub-tree

			// updating linked list
			next := n.Next
			prev := n.Prev
			if prev != nil {
				prev.Next = next
			}
			if next != nil {
				next.Prev = prev
			}
			n.Next = nil
			n.Prev = nil

			if t.maxC == n {
				t.maxC = prev
			}
			if t.minC == n {
				t.minC = next
			}

			return nil
		}

		if !t.isRed(n.right) && !t.isRed(n.right.left) {
			n = t.moveRedRight(n)
		}
		// h.right or one of its children red make

		if n.Key == key {
			// search hit, replacing the node with a successor
			rightMin := t.min(n.right)
			n.Key = rightMin.Key
			n.Value = rightMin.Value
			n.right = t.deleteMin(n.right)

			// global min will be updated automatically if required,
			// as we copy values from successor
		} else {
			if n.right == nil {
				// search miss
				return nil
			}

			n.right = t.delete(n.right, key)
		}
	}

	// balance
	if t.isRed(n.right) {
		n = t.rotateLeft(n)
	}
	if t.isRed(n.left) && t.isRed(n.left.left) {
		n = t.rotateRight(n)
	}
	if t.isRed(n.left) && t.isRed(n.right) {
		t.flipColors(n)
	}

	n.size = t.size(n.left) + 1 + t.size(n.right)
	return n
}

func (t *Tree) Keys(lo, hi float64) []float64 {
	if lo < t.Min() || hi > t.Max() {
		log.Panicln("keys out of range")
	}

	return t.keys(t.root, lo, hi)
}

func (t *Tree) keys(n *Node, lo, hi float64) []float64 {
	if n == nil {
		return nil
	}

	if n.Key < lo {
		return t.keys(n.right, lo, hi)
	} else if n.Key > hi {
		return t.keys(n.left, lo, hi)
	}

	l := t.keys(n.left, lo, hi)
	r := t.keys(n.right, lo, hi)

	keys := make([]float64, 0)
	if l != nil {
		keys = append(keys, l...)
	}
	keys = append(keys, n.Key)
	if r != nil {
		keys = append(keys, r...)
	}

	return keys
}

func (t *Tree) Print() {
	fmt.Println()
	t.print(t.root)
	fmt.Println()
}

func (t *Tree) print(n *Node) {
	if n == nil {
		return
	}

	if n.isRed {
		fmt.Printf("*")
	}
	fmt.Printf("%0.8f ", n.Key)

	t.print(n.left)
	t.print(n.right)
}
