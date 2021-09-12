package spider

import "time"

const (
	UserAgent = "OhYee-Spider"
	Timeout   = 120 * time.Second

	year2000 = 946656000 // 2000-01-01 00:00:00
)

var (
	linkKeys = []string{"link", "href", "url"}
)
