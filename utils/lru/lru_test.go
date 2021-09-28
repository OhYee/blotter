package lru

import (
	"testing"
)

func TestLRU(t *testing.T) {
	lru := NewLRU(5)

	assert := func(want, add string) {
		poped := lru.Add(add)
		if want != poped {
			t.Errorf("add %s, want %s, got %s", add, want, poped)
		}
	}

	assert("", "key1")
	assert("", "key2")
	assert("", "key3")
	assert("", "key4")
	assert("", "key5")
	assert("key1", "key6")
	lru.Visit("key2")
	assert("", "key3")
	assert("key4", "key7")
}
