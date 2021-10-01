package lru

import (
	"math/rand"
	"time"
)

type Map struct {
	m map[string]interface{}
	e *Heap
	l *LRU
}

func NewMap() *Map {
	return &Map{
		m: make(map[string]interface{}),
		e: nil,
		l: nil,
	}
}

func (m *Map) WithLRU(cap int) *Map {
	m.l = NewLRU(cap)
	return m
}

func (m *Map) WithExpired() *Map {
	m.e = NewHeap()
	return m
}

func (m *Map) Put(key string, value interface{}) {
	if m.l != nil {
		poped := m.l.Push(key)
		if poped != "" {
			m.remove(poped)
		}
	}
	m.m[key] = value
}

func (m *Map) PutWithExpired(key string, value interface{}, duration time.Duration) {
	m.Put(key, value)
	if m.e != nil {
		m.e.Push(key, time.Now().Add(duration).Unix())
	}
}

func (m *Map) Get(key string) (value interface{}, exists bool) {
	m.removeExpired()
	value, exists = m.m[key]
	return
}

func (m *Map) Len() int {
	m.removeExpired()
	return len(m.m)
}

func (m *Map) Keys() []string {
	m.removeExpired()
	keys := make([]string, 0, len(m.m))
	for k := range m.m {
		keys = append(keys, k)
	}
	return keys
}

func (m *Map) Delete(key string) {
	m.removeExpired()
	m.remove(key)
}

func (m *Map) gc() {
	// 10% 概率执行一次 gc
	if rand.Float64() > 0.1 {
		return
	}

	temp := m.m
	m.m = make(map[string]interface{})
	for k, v := range temp {
		m.m[k] = v
	}
}

func (m *Map) remove(key string) {
	if m.l != nil {
		m.l.Remove(key)
	}
	if m.e != nil {
		m.e.Remove(key)
	}
	delete(m.m, key)
	m.gc()
}

func (m *Map) removeExpired() {
	if m.e != nil {
		now := time.Now()
		keys := m.e.PopUntil(now.Unix())
		for _, k := range keys {
			m.remove(k)
		}
	}
}
