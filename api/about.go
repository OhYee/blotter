package api

import (
	"github.com/OhYee/blotter/api/pkg/variable"
	"github.com/OhYee/blotter/register"
)

// Experience for about api
type Experience struct {
	Name  string `json:"name"`
	Major string `json:"major"`
	Time  string `json:"time"`
}

// Award for about api
type Award struct {
	Name  string `json:"name"`
	Level string `json:"level"`
	Count int64  `json:"count"`
}

// AboutResponse response of avatar api
type AboutResponse struct {
	QQ          string       `json:"qq"`
	Github      string       `json:"github"`
	Email       string       `json:"email"`
	Zhihu       string       `json:"zhihu"`
	Author      string       `json:"author"`
	Quote       string       `json:"quote"`
	Description string       `json:"description"`
	Edu         []Experience `json:"edu"`
	Awards      []Award      `json:"awards"`
}

// About get avatar of emial
func About(context *register.HandleContext) (err error) {
	res := new(AboutResponse)

	data, err := variable.Get("github", "qq", "email", "zhihu", "author", "quote", "description", "edu", "awards")
	if err != nil {
		return
	}
	if err = data.SetString("qq", &res.QQ); err != nil {
		return
	}
	if err = data.SetString("github", &res.Github); err != nil {
		return
	}
	if err = data.SetString("email", &res.Email); err != nil {
		return
	}
	if err = data.SetString("zhihu", &res.Zhihu); err != nil {
		return
	}
	if err = data.SetString("author", &res.Author); err != nil {
		return
	}
	if err = data.SetString("quote", &res.Quote); err != nil {
		return
	}
	if err = data.SetString("description", &res.Description); err != nil {
		return
	}
	if err = data.SetArray("edu", &res.Edu); err != nil {
		return
	}
	if err = data.SetArray("awards", &res.Awards); err != nil {
		return
	}

	err = context.ReturnJSON(res)
	return
}
