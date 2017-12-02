package main

import (
	"container/heap"
)

type expiry struct {
	key string
	ttl int64 // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

func newExpiry(k string, ttl int64) *expiry {
	return &expiry{key: k, ttl: ttl}
}

// A expiryPriorityQueue implements heap.Interface and holds Items.
type expiryPriorityQueue []*expiry

func (pq expiryPriorityQueue) Len() int { return len(pq) }

func (pq expiryPriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].ttl < pq[j].ttl
}

func (pq expiryPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *expiryPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*expiry)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *expiryPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *expiryPriorityQueue) add(e *expiry) {
	heap.Push(pq, e)
}

// update modifies the ttl of an expiry in the queue.
func (pq *expiryPriorityQueue) update(e *expiry, ttl int64) {
	e.ttl = ttl
	heap.Fix(pq, e.index)
}

func (pq *expiryPriorityQueue) del(e *expiry) {
	heap.Remove(pq, e.index)
}
