package spider

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/OhYee/blotter/output"
	"github.com/chromedp/chromedp"
)

var client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: false,
	},
	Timeout: Timeout,
}

func getHTML(u string) (content string, isJSON bool) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		output.ErrOutput.Println(u, err)
		return
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Add("Accept-Encoding", "identity")
	// https://dpjeep.com/2019/06/10/golangzhi-http-eofxiang-jie/
	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		output.ErrOutput.Println(u, err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		output.ErrOutput.Println(u, err)
		return
	}

	content = string(b)
	isJSON = strings.Contains(resp.Header.Get("Content-Type"), "json")
	return
}

func getChromePath() string {
	env := os.Environ()
	for _, s := range env {
		ss := strings.Split(s, "=")
		if len(ss) >= 2 {
			key := ss[0]
			value := strings.Join(ss[1:], "")
			if strings.ToUpper(key) == "CHROME_PATH" {
				return value
			}
		}
	}
	return ""
}
func getHTMLWithJS(u string) string {
	opts := []func(*chromedp.ExecAllocator){
		chromedp.Flag("headless", true),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("default-browser-check ", false),
		chromedp.UserAgent(UserAgent),
	}

	chromePath := getChromePath()
	if len(chromePath) > 0 {
		opts = append(
			opts,
			chromedp.ExecPath(
				fmt.Sprintf(
					"%s/%s",
					strings.TrimRight(chromePath, "/"),
					"chrome",
				),
			),
		)
	}

	ctx, cancel := chromedp.NewExecAllocator(
		context.Background(),

		opts...,
	)
	defer cancel()

	ctx2, cancel2 := chromedp.NewContext(
		ctx,
	)
	defer cancel2()

	ctx3, cancel3 := context.WithTimeout(ctx2, Timeout)
	defer cancel3()

	var res string
	err := chromedp.Run(
		ctx3,
		chromedp.Navigate(u),
		chromedp.OuterHTML("html", &res, chromedp.ByQuery),
	)
	if err != nil {
		output.ErrOutput.Println(u, err)
	}
	return res
}
