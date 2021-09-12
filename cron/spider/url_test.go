package spider

import (
	"fmt"
	"testing"
)

func Test_makeAbsURL(t *testing.T) {
	tests := []struct {
		host string
		path string
		want string
	}{
		{
			host: "https://www.oyohyee.com",
			path: "https://www.oyohyee.com/index.html",
			want: "https://www.oyohyee.com/index.html",
		},
		{
			host: "https://www.oyohyee.com:8080",
			path: "https://www.oyohyee.com/index.html",
			want: "https://www.oyohyee.com/index.html",
		},
		{
			host: "https://www.oyohyee.com",
			path: "/index.html",
			want: "https://www.oyohyee.com/index.html",
		},
		{
			host: "http://www.oyohyee.com/",
			path: "index.html",
			want: "http://www.oyohyee.com/index.html",
		},
		{
			host: "",
			path: "index.html",
			want: "index.html",
		},
		{
			host: "https://www.ohyee.cc/index.html",
			path: "!@#$%",
			want: "https://www.ohyee.cc/index.html",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s | %s", tt.host, tt.path), func(t *testing.T) {
			if got := makeAbsURL(tt.host, tt.path); got != tt.want {
				t.Errorf("getHost() = %v, want %v", got, tt.want)
			}
		})
	}
}
