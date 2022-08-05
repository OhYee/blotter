package main

import (
	"fmt"

	"github.com/OhYee/blotter/utils/geoip"
)

func main() {
	fmt.Println(geoip.GetPositionFromIP("1.2.3.4"))
}
