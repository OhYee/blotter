package geoip

import (
	_ "embed"
	"net/http"
	"testing"
)

func TestGetIPFromHeader(t *testing.T) {
	tests := []struct {
		name   string
		header *http.Header
		want   string
	}{
		{
			name:   "empty header",
			header: &http.Header{},
			want:   "",
		},
		{
			name: "empty X-Forwarded-For",
			header: &http.Header{
				"X-Forwarded-For": []string{},
			},
			want: "",
		},
		{
			name: "one ip",
			header: &http.Header{
				"X-Forwarded-For": []string{"127.0.0.1"},
			},
			want: "127.0.0.1",
		},
		{
			name: "one ip",
			header: &http.Header{
				"X-Forwarded-For": []string{"127.0.0.1"},
			},
			want: "127.0.0.1",
		},
		{
			name: "two ip in arr",
			header: &http.Header{
				"X-Forwarded-For": []string{"127.0.0.1", "127.0.0.2"},
			},
			want: "127.0.0.1",
		},
		{
			name: "two ip in string",
			header: &http.Header{
				"X-Forwarded-For": []string{"127.0.0.1, 127.0.0.2"},
			},
			want: "127.0.0.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetIPFromHeader(tt.header); got != tt.want {
				t.Errorf("GetIPFromHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
