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
}
