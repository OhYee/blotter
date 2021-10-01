package lru

import (
	"sync/atomic"
)

type LRU struct {
	heap    *Heap
	counter int64
	cap     int
}

func NewLRU(cap int) *LRU {
	return &LRU{
		counter: 0,
		heap:    NewHeap(),
		cap:     cap,
	}
}

func (l *LRU) Push(key string) string {
	id := atomic.AddInt64(&l.counter, 1)

	poped := ""
	for l.heap.Len() >= l.cap && !l.Has(key) {
		poped = l.Pop()
	}

	l.heap.Push(key, id)
	return poped
}

func (l *LRU) Pop() string {
	return l.heap.Pop()
}

func (l *LRU) Has(key string) bool {
	return l.heap.Has(key)
}

func (l *LRU) Remove(key string) {
	l.heap.Remove(key)
}
