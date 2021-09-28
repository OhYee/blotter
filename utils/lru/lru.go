package lru

import (
	"container/heap"
	"sync/atomic"
)

type LRU struct {
	counter int64
	heap    *keyValueHeap
	m       map[string]*keyValue
	cap     int
}

func NewLRU(cap int) *LRU {
	return &LRU{
		counter: 0,
		heap:    newKeyValueHeap(0),
		m:       make(map[string]*keyValue),
		cap:     cap,
	}
}

func (l *LRU) Add(key string) string {
	id := atomic.AddInt64(&l.counter, 1)

	if item, exists := l.m[key]; exists {
		item.expiredTime = id
		heap.Fix(l.heap, item.pos)
		return ""
	}

	poped := ""
	if l.heap.Len() == l.cap {
		poped = heap.Pop(l.heap).(*keyValue).key
		delete(l.m, poped)
	}

	item := &keyValue{
		key:         key,
		expiredTime: id,
	}
	l.m[key] = item
	heap.Push(l.heap, item)
	return poped
}

func (l *LRU) Visit(key string) bool {
	item, exist := l.m[key]
	if exist {
		item.expiredTime = atomic.AddInt64(&l.counter, 1)
		heap.Fix(l.heap, item.pos)
	}
	return exist
}
