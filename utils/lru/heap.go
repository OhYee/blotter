package lru

import (
	"container/heap"
)

type keyValue struct {
	key   string
	value int64
	pos   int
}

type keyValueHeap []*keyValue

var _ heap.Interface = (*keyValueHeap)(nil)

func newKeyValueHeap(l int) *keyValueHeap {
	e := keyValueHeap(make([]*keyValue, 0, l))
	return &e
}
func (h keyValueHeap) At(i int) int64 { return h[i].value }
func (h keyValueHeap) Len() int       { return len(h) }
func (h keyValueHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].pos = i
	h[j].pos = j
}
func (h keyValueHeap) Less(i, j int) bool     { return h[i].value < h[j].value }
func (h *keyValueHeap) Push(item interface{}) { *h = append(*h, item.(*keyValue)) }
func (h *keyValueHeap) Pop() interface{} {
	l := h.Len()
	t := (*h)[l-1]
	(*h) = (*h)[:l-1]
	return t
}

type Heap struct {
	heap *keyValueHeap
	m    map[string]*keyValue
}

func NewHeap() *Heap {
	return &Heap{
		heap: newKeyValueHeap(0),
		m:    make(map[string]*keyValue),
	}
}
func (h *Heap) Len() int {
	return len(h.m)
}

func (h *Heap) Push(key string, value int64) {
	if item, exists := h.m[key]; exists {
		item.value = value
		heap.Fix(h.heap, item.pos)
		return
	}

	item := &keyValue{
		key:   key,
		value: value,
	}
	h.m[key] = item
	heap.Push(h.heap, item)
}

func (h *Heap) Pop() string {
	if h.heap.Len() == 0 {
		return ""
	}
	item := heap.Pop(h.heap).(*keyValue)
	delete(h.m, item.key)
	h.gc()
	return item.key
}

func (h *Heap) Has(key string) bool {
	_, exists := h.m[key]
	return exists
}

func (h *Heap) Remove(key string) {
	item, exist := h.m[key]
	if exist {
		heap.Remove(h.heap, item.pos)
		delete(h.m, key)
	}
}

func (h *Heap) PopUntil(value int64) []string {
	keys := make([]string, 0)
	for h.Len() > 0 && h.heap.At(0) <= value {
		keys = append(keys, h.Pop())
	}
	return keys
}

func (h *Heap) gc() {
	temp := h.m
	h.m = make(map[string]*keyValue)
	for k, v := range temp {
		h.m[k] = v
	}
}
