package lru

import (
	"container/heap"
)

var _ heap.Interface = (*keyValueHeap)(nil)

type Expired struct {
	heap *keyValueHeap
}

func NewExpired() *Expired {
	return &Expired{
		heap: newKeyValueHeap(0),
	}
}

func (e *Expired) Add(key string, expiredTime int64) {
	heap.Push(
		e.heap,
		&keyValue{
			key:         key,
			expiredTime: expiredTime,
		},
	)
}

func (e *Expired) Expire(now int64) []string {
	keys := make([]string, 0)
	for e.heap.Len() > 0 && (*e.heap)[0].expiredTime <= now {
		keys = append(keys, heap.Pop(e.heap).(*keyValue).key)
	}
	return keys
}
