package serverchan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/OhYee/blotter/api/pkg/variable"
	"github.com/OhYee/blotter/output"
)

// Notify send a notification to ServerChan
func Notify(title string, message string) {
	v, err := variable.Get("server_chan")
	if err != nil {
		output.Err(err)
		return
	}
	token, ok := v.GetString("server_chan")
	if token == "" || !ok {
		output.ErrOutput.Println("Can not get token for ServerChan")
		return
	}
	params := url.Values{}
	params.Set("text", title)
	params.Set("desp", message)

	resp, err := http.PostForm(
		fmt.Sprintf("http://sc.ftqq.com/%s.send", token),
		params,
	)
	if err != nil {
		output.Err(err)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		output.Err(err)
		return
	}
	res := make(map[string]interface{})
	err = json.Unmarshal(b, &res)
	if err != nil {
		output.Err(err)
		return
	}
	output.Debug("%+b", res)
	return
}
