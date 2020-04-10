package api

import "github.com/OhYee/blotter/register"

// Robots robots.txt
func Robots(context register.HandleContext) (err error) {
	robots := `User-agent: *
Disallow: 
Sitemap: /sitemap.xml`
	context.ReturnText(robots)
	return
}
