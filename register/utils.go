package register

import (
	"net/http"
	"strings"
)

const (
	X_Real_IP       = "X-Real-Ip"
	X_FORWARDED_FOR = "X-Forwarded-For"
)

func getIPFromHeader(header *http.Header, headerName string) string {
	IPs := header.Values(headerName)
	if len(IPs) > 0 {
		remoteIP := IPs[0]
		arr := strings.Split(remoteIP, ",")
		if len(arr) > 0 {
			return strings.TrimSpace(arr[0])
		}
	}
	return ""
}
