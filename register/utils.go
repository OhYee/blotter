package register

import (
	"net/http"
	"strings"
)

const (
	_X_Real_IP       = "X-Real-Ip"
	_X_FORWARDED_FOR = "X-Forwarded-For"
)

var (
	ipHeaders = []string{_X_FORWARDED_FOR, _X_Real_IP}
)

func getIPFromHeader(header *http.Header, headerName string) string {
	IPs := header.Values(headerName)

	if len(IPs) > 0 {
		remoteIP := IPs[len(IPs)-1]
		arr := strings.Split(remoteIP, ",")
		if len(arr) > 0 {
			return strings.TrimSpace(arr[0])
		}
	}
	return ""
}
