package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/OhYee/blotter/api/pkg/variable"
	"github.com/OhYee/blotter/register"
	"github.com/OhYee/rainbow/errors"
	mv2ray "github.com/v2fly/vmessping/miniv2ray"
)

func githubRepo(username string, page int) (res []map[string]interface{}, err error) {
	var client *http.Client = nil

	vars, err := variable.Get("vmess")
	vmessLink, _ := vars.GetString("vmess")
	fmt.Println(vmessLink)
	if vmessLink != "" {
		server, err := mv2ray.StartV2Ray(vmessLink, false, true)
		if err == nil {
			err = server.Start()
			if err == nil {
				defer server.Close()
				client, _ = mv2ray.CoreHTTPClient(server, time.Duration(10*time.Second))
			}
		}
	}

	if client == nil {
		fmt.Println("use default")
		client = http.DefaultClient
	}

	resp, err := client.Get(fmt.Sprintf("https://proxy.ohyee.cc/api.github.com/users/%s/repos?page=%d", username, page))
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("Status code %d: %s", resp.StatusCode, resp.Body)
		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	res = make([]map[string]interface{}, 0)
	err = json.Unmarshal(bytes, &res)
	return
}

// GithubReposRequest request for GithubRepos api
type GithubReposRequest struct {
	Username string `json:"username"`
}

// GithubReposResponse response for GithubRepos api
type GithubReposResponse struct {
	Repos []map[string]interface{} `json:"repos"`
}

// GithubRepos get github repos of {username}
func GithubRepos(context register.HandleContext) (err error) {
	args := new(GithubReposRequest)
	res := new(GithubReposResponse)
	context.RequestArgs(args)

	res.Repos = make([]map[string]interface{}, 0)

	page := 1
	for {
		var repos []map[string]interface{}
		if repos, err = githubRepo(args.Username, page); err != nil {
			return
		}
		if len(repos) == 0 {
			break
		}

		res.Repos = append(res.Repos, repos...)
		page++
	}

	err = context.ReturnJSON(res)
	return
}
