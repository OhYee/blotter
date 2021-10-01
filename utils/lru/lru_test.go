package lru

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestLRU(t *testing.T) {
	lru := NewLRU(5)

	assertPush := func(want, add string) {
		poped := lru.Push(add)
		assert.Equal(t, want, poped)
	}

	assertPush("", "key1")
	assertPush("", "key2")
	assertPush("", "key3")
	assertPush("", "key4")
	assertPush("", "key5")
	assertPush("key1", "key6")
	lru.Push("key2")
	assertPush("", "key3")
	assertPush("key4", "key7")
	lru.Remove("key5")
	assertPush("", "key8")

	assert.Equal(t, "key6", lru.Pop())
	assert.Equal(t, "key2", lru.Pop())
	assert.Equal(t, "key3", lru.Pop())
	assert.Equal(t, "key7", lru.Pop())
	assert.Equal(t, "key8", lru.Pop())
	assert.Equal(t, "", lru.Pop())
}
