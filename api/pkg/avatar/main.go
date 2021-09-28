package avatar

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const defaultAvatar = "https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png"

func makeHTTPClient() *http.Client {
	keys := []string{"HTTP_PROXY", "HTTPS_PROXY", "PROXY", "http_proxy", "https_proxy", "proxy"}
	proxy := ""
	for _, key := range keys {
		proxy = os.Getenv(key)
		if proxy != "" {
			break
		}
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	if proxy != "" {
		client.Transport = &http.Transport{
			Proxy: func(_ *http.Request) (*url.URL, error) {
				return url.Parse(proxy)
			},
		}
	}

	return client
}

type githubSearchAPI struct {
	TotalCount int              `json:"total_count"`
	Incomplete bool             `json:"incomplete_results"`
	Items      []githubUserInfo `json:"items"`
}

type githubUserInfo struct {
	Login             string  `json:"login"`
	ID                int     `json:"id"`
	NodeID            string  `json:"node_id"`
	AvatarURL         string  `json:"avatar_url"`
	GravatarID        string  `json:"gravatar_id"`
	URL               string  `json:"url"`
	HTMLlURL          string  `json:"html_url"`
	FollowersURL      string  `json:"followers_url"`
	FollowingURL      string  `json:"following_url"`
	GistsURL          string  `json:"gists_url"`
	StarredURL        string  `json:"starred_url"`
	SubscriptionsURL  string  `json:"subscriptions_url"`
	OrganizationsURL  string  `json:"organizations_url"`
	ReposURL          string  `json:"repos_url"`
	EventsURL         string  `json:"events_url"`
	ReceivedEventsURL string  `json:"received_events_url"`
	Type              string  `json:"type"`
	SiteAdmin         bool    `json:"site_admin"`
	Score             float64 `json:"score"`
}

// GetGithubAvatar get the github avatar of the email
func GetGithubAvatar(email string) (avatar string) {
	rep, err := makeHTTPClient().Get(fmt.Sprintf("https://api.github.com/search/users?q=%s", email))
	if err != nil {
		return
	}
	defer rep.Body.Close()

	data, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return
	}

	res := githubSearchAPI{}

	json.Unmarshal(data, &res)
	if res.TotalCount > 0 {
		avatar = res.Items[0].AvatarURL
	}

	return
}

// GetGravatar get Gravatar of the gmail
func GetGravatar(email string) (avatar string) {
	hash := md5.New()
	hash.Write([]byte(strings.ToLower(email)))

	avatar = fmt.Sprintf(
		"https://www.gravatar.com/avatar/%s?size=640&default=%s",
		hex.EncodeToString(hash.Sum(nil)),
		url.QueryEscape(defaultAvatar),
	)

	req, err := makeHTTPClient().Get(avatar)
	if err != nil || req.StatusCode == 302 {
		avatar = ""
	}

	return
}

// Get avatar url of the email
func Get(email string) (avatar string) {
	avatar = GetGithubAvatar(email)
	if avatar == "" {
		avatar = GetGravatar(email)
	}
	if avatar == "" {
		avatar = defaultAvatar
	}
	return
}
