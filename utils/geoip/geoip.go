package geoip

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/OhYee/rainbow/log"
)

type IPInfo struct {
	Status      string
	Country     string
	CountryCode string
	Region      string
	RegionName  string
	City        string
	Zip         string
	Lat         float64
	Lon         float64
	TimeZone    string
	Isp         string
	Org         string
	As          string
	AsName      string
	Query       string
	Offset      int
	Currency    string
	Mobile      bool
	Proxy       bool
	Hosting     bool
}

func GetPositionFromIP(ip string) string {
	resp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s?fields=66842623&lang=zh-CN", ip))
	if err != nil {
		log.Error.Printf("failed to get response from ip-api, due to %s", err)
		return ""
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error.Printf("failed to read response from ip-api, due to %s", err)
		return ""
	}

	ipInfo := new(IPInfo)
	err = json.Unmarshal(b, ipInfo)
	if err != nil {
		log.Error.Printf("failed to read json from ip-api response body, due to %s", err)
		return ""
	}

	positionArr := make([]string, 0, 3)
	if ipInfo.Country != "" {
		positionArr = append(positionArr, ipInfo.Country)
	}
	if ipInfo.RegionName != "" {
		positionArr = append(positionArr, ipInfo.RegionName)
	}
	if ipInfo.City != "" {
		positionArr = append(positionArr, ipInfo.City)
	}
	if ipInfo.Mobile {
		positionArr = append(positionArr, "移动端")
	}
	if ipInfo.Proxy {
		positionArr = append(positionArr, "代理")
	}

	return strings.Join(positionArr, ",")
}
