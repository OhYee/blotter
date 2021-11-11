package spider

import (
	"fmt"
	"testing"
)

func Test_toInt64(t *testing.T) {

	tests := []struct {
		args interface{}
		want int64
	}{
		{100, 100},
		{100.1, 100},
		{uint8(100), 100},
		{"a100", 0},
		{1.2, 1},
		{1e5, 100000},
		{1.2e5, 120000},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v -> %#v", tt.args, tt.want), func(t *testing.T) {
			if got := toInt64(tt.args); got != tt.want {
				t.Errorf("toInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}
