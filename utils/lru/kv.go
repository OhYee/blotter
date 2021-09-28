package lru

type keyValue struct {
	key         string
	expiredTime int64
	pos         int
}

type keyValueHeap []*keyValue

func newKeyValueHeap(l int) *keyValueHeap {
	e := keyValueHeap(make([]*keyValue, 0, l))
	return &e
}
func (h keyValueHeap) Len() int { return len(h) }
func (h keyValueHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].pos = i
	h[j].pos = j
}
func (h keyValueHeap) Less(i, j int) bool     { return h[i].expiredTime < h[j].expiredTime }
func (h *keyValueHeap) Push(item interface{}) { *h = append(*h, item.(*keyValue)) }
func (h *keyValueHeap) Pop() interface{} {
	l := h.Len()
	t := (*h)[l-1]
	(*h) = (*h)[:l-1]
	return t
}
