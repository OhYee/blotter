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
	"regexp"
	"strings"
	"time"

	"github.com/OhYee/blotter/utils/lru"
)

const DefaultAvatar = "https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png"
const DefaultAvatarMD5 = "1fbf1eeb622038a1ea2e62036d33788a"

func md5Encode(source []byte) string {
	hash := md5.New()
	hash.Write(source)
	return hex.EncodeToString(hash.Sum(nil))
}

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

// GetCravatar get Gravatar of the gmail
func GetCravatar(email string) (avatar string) {
	avatar = fmt.Sprintf(
		"https://cravatar.cn/avatar/%s?s=640&d=404",
		md5Encode([]byte(strings.ToLower(email))),
	)

	resp, err := makeHTTPClient().Get(avatar)
	if err != nil || resp == nil || resp.StatusCode == 404 {
		avatar = ""
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || md5Encode(body) == DefaultAvatarMD5 {
		avatar = ""
	}
	return
}

// GetGravatar get Gravatar of the gmail
func GetGravatar(email string) (avatar string) {
	avatar = fmt.Sprintf(
		"https://www.gravatar.com/avatar/%s?size=640&default=%s",
		md5Encode([]byte(strings.ToLower(email))),
		url.QueryEscape(DefaultAvatar),
	)

	req, err := makeHTTPClient().Get(avatar)
	if err != nil || req.StatusCode == 302 {
		avatar = ""
	}

	return
}

var qqRegexp = regexp.MustCompile(`^\d+$`)

func GetQQAvatar(email string) (avatar string) {
	if strings.HasSuffix(email, "@qq.com") {
		qq := email[:len(email)-7]
		if qqRegexp.MatchString(qq) {
			avatar = fmt.Sprintf("https://q1.qlogo.cn/g?b=qq&nk=%s&s=640", qq)
		}
	}
	return
}

type avatarGetter = func(string) string

var avatarGetters = []avatarGetter{
	GetGithubAvatar,
	GetCravatar,
	GetGravatar,
	GetQQAvatar,
}

var avatarMap = lru.NewMap().WithLRU(50).WithExpired()

// Get avatar url of the email
func Get(email string) (avatar string) {
	value, exists := avatarMap.Get(email)
	if exists {
		avatar = value.(string)
		return
	}

	for _, getter := range avatarGetters {
		avatar = getter(email)
		if avatar != "" {
			break
		}
	}
	if avatar == "" {
		avatar = DefaultAvatar
	}

	if avatar != "" {
		avatarMap.PutWithExpired(email, avatar, time.Hour*24*7)
	}
	return
}
