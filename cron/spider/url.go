package spider

import (
	"net/url"
)

// makeAbsURL returns the url concat by host and path
func makeAbsURL(host, path string) string {
	pathObj, err := url.Parse(path)
	if err != nil {
		return host
	}
	if pathObj.IsAbs() {
		return pathObj.String()
	}
	hostObj, err := url.Parse(host)
	if err == nil {
		pathObj.Scheme = hostObj.Scheme
		pathObj.Host = hostObj.Host
	}
	return pathObj.String()
}
