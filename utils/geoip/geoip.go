package geoip

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/OhYee/blotter/output"
)

type ipAPIResponse struct {
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

func ipAPI(ip string) string {
	apiName := "ip-api"

	resp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s?fields=66842623&lang=zh-CN", ip))
	if err != nil {
		output.ErrOutput.Printf("failed to get response from %s for %s due to %s", apiName, ip, err)
		return ""
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		output.ErrOutput.Printf("failed to read response from %s for %s due to %s", apiName, ip, err)
		return ""
	}

	ipInfo := new(ipAPIResponse)
	err = json.Unmarshal(b, ipInfo)
	if err != nil {
		output.ErrOutput.Printf("failed to read json from %s for %s, due to %s", apiName, ip, err)
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

type ohyeeAPIResponse struct {
	Success   bool
	IP        string
	Continent string
	Country   string
	Region    string
	City      string
	Mobile    bool
	Proxy     bool
	Lat       float64
	Lon       float64
}

func ohyeeAPI(ip string) string {
	const apiName = "ohyee api"

	resp, err := http.Get(fmt.Sprintf("http://fc.ohyee.cc/ip?ip=%s", ip))
	if err != nil {
		output.ErrOutput.Printf("failed to get response from %s for %s, due to %s", apiName, ip, err)
		return ""
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		output.ErrOutput.Printf("failed to read response from %s for %s, due to %s", apiName, ip, err)
		return ""
	}

	ipInfo := new(ohyeeAPIResponse)
	err = json.Unmarshal(b, ipInfo)
	if err != nil {
		output.ErrOutput.Printf("failed to read json from %s for %s, due to %s", apiName, ip, err)
		return ""
	}

	positionArr := make([]string, 0, 3)
	if ipInfo.Country != "" {
		positionArr = append(positionArr, ipInfo.Country)
	}
	if ipInfo.Region != "" {
		positionArr = append(positionArr, ipInfo.Region)
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

var apis = []func(string) string{
	ipAPI,
	ohyeeAPI,
}

func GetPositionFromIP(ip string) string {
	for _, f := range apis {
		if pos := f(ip); pos != "" {
			return pos
		}
	}

	return ""
}
