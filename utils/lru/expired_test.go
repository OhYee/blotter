package lru

import (
	"sort"
	"strings"
	"testing"
)

func TestExpired(t *testing.T) {
	type test struct {
		key  string
		time int64
		now  int64
		keys []string
	}

	tests := []test{
		{"key1", 50, 10, []string{}},
		{"key2", 30, 20, []string{}},
		{"key3", 50, 30, []string{"key2"}},
		{"key4", 20, 40, []string{"key4"}},
		{"key5", 10, 50, []string{"key1", "key3", "key5"}},
		{"key6", 50, 60, []string{"key6"}},
		{"key7", 100, 70, []string{}},
	}

	expired := NewExpired()
	for idx, tt := range tests {
		expired.Add(tt.key, tt.time)
		keys := expired.Expire(tt.now)

		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
		sort.Slice(tt.keys, func(i, j int) bool { return tt.keys[i] < tt.keys[j] })
		if strings.Join(keys, "|") != strings.Join(tt.keys, "|") {
			t.Errorf("%d | keys got %v, want %v", idx, keys, tt.keys)
		}
	}
}
