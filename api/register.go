package api

import (
	"github.com/OhYee/blotter/register"
)

// Register api
func Register() {
	register.Register(
		"friends",
		Friends,
	)
	register.Register(
		"menus",
		Menus,
	)
	register.Register(
		"post",
		Post,
	)
	register.Register(
		"admin/post",
		PostAdmin,
	)
	register.Register(
		"admin/post/edit",
		PostEdit,
	)
	register.Register(
		"post/existed",
		PostExisted,
	)
	register.Register(
		"posts",
		Posts,
	)
	register.Register(
		"admin/posts",
		PostsAdmin,
	)
	register.Register(
		"admin/post/delete",
		PostDelete,
	)
	register.Register(
		"markdown",
		Markdown,
	)
	register.Register(
		"markdown/ws",
		MarkdownWS,
	)
	register.Register(
		"comments",
		Comments,
	)
	register.Register(
		"layout",
		Layout,
	)
	register.Register(
		"tags",
		Tags,
	)
	register.Register(
		"avatar",
		Avatar,
	)
	register.Register(
		"comment/add",
		CommentAdd,
	)
	register.Register(
		"login",
		Login,
	)
	register.Register(
		"logout",
		Logout,
	)
	register.Register(
		"info",
		Info,
	)
	register.Register(
		"admin/tag/edit",
		TagEdit,
	)
	register.Register(
		"admin/tag/delete",
		TagDelete,
	)
	register.Register(
		"tag/existed",
		TagExisted,
	)
	register.Register(
		"tag",
		Tag,
	)
	register.Register(
		"robots.txt",
		Robots,
	)
	register.Register(
		"sitemap.txt",
		SitemapTXT,
	)
	register.Register(
		"sitemap.xml",
		SitemapXML,
	)
	register.Register(
		"rss.xml",
		RSSXML,
	)
	register.Register(
		"admin/friends/set",
		SetFriends,
	)
	register.Register(
		"view",
		View,
	)
	register.Register(
		"admin/menus/set",
		SetMenus,
	)
	register.Register(
		"about",
		About,
	)
	register.Register(
		"admin/variables",
		Variables,
	)
	register.Register(
		"admin/variables/set",
		VariablesSet,
	)
	register.Register(
		"admin/comments",
		AdminComments,
	)
	register.Register(
		"admin/comment/set",
		AdminCommentSet,
	)
	register.Register(
		"admin/comment/delete",
		AdminCommentDelete,
	)
	register.Register(
		"user/qq_connect",
		QQ,
	)
	register.Register(
		"user/jump_to_qq",
		JumpToQQ,
	)
	register.Register(
		"user/set",
		SetUser,
	)
	register.Register(
		"user/username",
		CheckUsername,
	)
	register.Register(
		"user/register",
		RegisterUser,
	)
	register.Register(
		"user/qq_avatar",
		SyncQQAvatar,
	)
	register.Register(
		"notification/ws",
		WebSocket,
	)
	register.Register(
		"user/github_connect",
		Github,
	)
	register.Register(
		"user/jump_to_github",
		JumpToGithub,
	)
	register.Register(
		"users",
		Users,
	)
	register.Register(
		"admin/user/reset_password",
		ResetPassword,
	)
	register.Register(
		"github/repos",
		GithubRepos,
	)
	register.Register(
		"travels",
		TravelsGet,
	)
	register.Register(
		"travels/url",
		TravelsGetByURL,
	)
	register.Register(
		"travels/set",
		TravelsSet,
	)
	register.Register(
		"qiniu/buckets",
		GetBuckets,
	)
	register.Register(
		"qiniu/images",
		GetImages,
	)
	register.Register(
		"qiniu/token",
		GetQiniuToken,
	)
	register.Register(
		"qiniu/image/delete",
		DeleteImage,
	)
	register.Register(
		"qiniu/image/rename",
		RenameImage,
	)
	register.Register(
		"wechat",
		WechatCheckPermission,
	)
}
