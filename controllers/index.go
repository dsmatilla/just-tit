package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"html/template"
)


// JTCache Just-tit cache
var JTCache cache.Cache

// JTVideo Just-tit video struct
type JTVideo struct {
	ID		     	string
	Domain	     	template.URL
	Title        	string
	Description		string
	Thumb		 	string
	Thumbs			[]string
	Embed 			template.HTML
	URL				template.URL
	Provider		string
	Rating			string
	Ratings			string
	Duration		string
	Views			string
	Width			string
	Height			string
	Segment			string
	PublishDate		string
	Type 			string
	Tags			[]string
	Categories		[]string
	Pornstars		[]string
	ExternalURL		string
	ExternalID		string
}

// IndexController Beego Controller
type IndexController struct {
	beego.Controller
}

// Get Index Controller
func (c *IndexController) Get() {
	c.Data["PageTitle"] = "Just Tit"
	c.Data["PageMetaDesc"] = "The most optimized adult video search engine"

	search := c.GetString("s")
	if len(search) > 0 {
		c.Redirect(search + ".html", 301)
	}

	//c.Layout = "index.tpl"
	c.TplName = "index.tpl"
}
