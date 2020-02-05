package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
)

var JTCache cache.Cache

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	c.Data["PageTitle"] = "Just Tit"
	c.Data["PageMetaDesc"] = "The most optimized adult video search engine"

	search := c.GetString("s")
	if len(search) > 0 {
		c.Redirect(search + ".html", 301)
	}

	c.TplName = "index.html"
}
