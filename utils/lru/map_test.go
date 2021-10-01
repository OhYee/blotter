package lru

import (
	"testing"
	"time"

	"github.com/alecthomas/assert"
)

func doMap(m *Map) {
	m.Put("key1", "value1")
	m.Put("key2", "value2")
	m.PutWithExpired("key3", "value3", time.Millisecond*200)
	m.Put("key4", "value4")
	m.Put("key5", "value5")
}

func TestMap(t *testing.T) {
	m := NewMap()
	doMap(m)
	var value interface{}
	var exists bool

	value, exists = m.Get("key1")
	assert.Equal(t, "value1", value)
	assert.True(t, exists)

	value, exists = m.Get("key2")
	assert.Equal(t, "value2", value)
	assert.True(t, exists)

	value, exists = m.Get("key3")
	assert.Equal(t, "value3", value)
	assert.True(t, exists)

	value, exists = m.Get("key4")
	assert.Equal(t, "value4", value)
	assert.True(t, exists)

	value, exists = m.Get("key5")
	assert.Equal(t, "value5", value)
	assert.True(t, exists)

	keys := m.Keys()
	assert.Equal(t, 5, m.Len())
	assert.Equal(t, 5, len(keys))
	for _, key := range keys {
		m.Delete(key)
	}
	assert.Equal(t, 0, m.Len())
}

func TestMapLRU(t *testing.T) {
	m := NewMap().WithLRU(3)
	doMap(m)
	var value interface{}
	var exists bool

	value, exists = m.Get("key1")
	assert.Nil(t, value)
	assert.False(t, exists)

	value, exists = m.Get("key2")
	assert.Nil(t, value)
	assert.False(t, exists)

	value, exists = m.Get("key3")
	assert.Equal(t, "value3", value)
	assert.True(t, exists)

	value, exists = m.Get("key4")
	assert.Equal(t, "value4", value)
	assert.True(t, exists)

	value, exists = m.Get("key5")
	assert.Equal(t, "value5", value)
	assert.True(t, exists)

	keys := m.Keys()
	assert.Equal(t, 3, m.Len())
	assert.Equal(t, 3, len(keys))
	for _, key := range keys {
		m.Delete(key)
	}
	assert.Equal(t, 0, m.Len())
}

func TestMapExpired(t *testing.T) {
	m := NewMap().WithExpired()
	doMap(m)
	time.Sleep(500 * time.Millisecond)

	var value interface{}
	var exists bool

	value, exists = m.Get("key1")
	assert.Equal(t, "value1", value)
	assert.True(t, exists)

	value, exists = m.Get("key2")
	assert.Equal(t, "value2", value)
	assert.True(t, exists)

	value, exists = m.Get("key3")
	assert.Nil(t, value)
	assert.False(t, exists)

	value, exists = m.Get("key4")
	assert.Equal(t, "value4", value)
	assert.True(t, exists)

	value, exists = m.Get("key5")
	assert.Equal(t, "value5", value)
	assert.True(t, exists)

	keys := m.Keys()
	assert.Equal(t, 4, m.Len())
	assert.Equal(t, 4, len(keys))
	for _, key := range keys {
		m.Delete(key)
	}
	assert.Equal(t, 0, m.Len())
}

func TestMapLRUExpired(t *testing.T) {
	m := NewMap().WithLRU(3).WithExpired()
	doMap(m)
	time.Sleep(500 * time.Millisecond)

	var value interface{}
	var exists bool

	value, exists = m.Get("key1")
	assert.Nil(t, value)
	assert.False(t, exists)

	value, exists = m.Get("key2")
	assert.Nil(t, value)
	assert.False(t, exists)

	value, exists = m.Get("key3")
	assert.Nil(t, value)
	assert.False(t, exists)

	value, exists = m.Get("key4")
	assert.Equal(t, "value4", value)
	assert.True(t, exists)

	value, exists = m.Get("key5")
	assert.Equal(t, "value5", value)
	assert.True(t, exists)

	keys := m.Keys()
	assert.Equal(t, 2, m.Len())
	assert.Equal(t, 2, len(keys))
	for _, key := range keys {
		m.Delete(key)
	}
	assert.Equal(t, 0, m.Len())
}
