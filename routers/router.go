package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/client/cache"
	_ "github.com/beego/beego/v2/client/cache/redis"
	"just-tit/controllers"
	"os"
)

func init() {
	// Initialize memory cache
	redisHost := os.Getenv("REDISHOST")
	redisName := os.Getenv("REDISNAME")
	redisDBNum := os.Getenv("REDISDBNUM")

	if redisHost != "" {
		controllers.JTCache, _ = cache.NewCache("redis", `{"key":"`+redisName+`","conn":"`+redisHost+`","dbNum":"`+redisDBNum+`"}`)
	} else {
		controllers.JTCache, _ = cache.NewCache("memory", `{"interval":60}`)
	}

	beego.Router("/", &controllers.IndexController{})
	beego.Router("/*.html", &controllers.SearchController{})
	beego.Router("/images/*", &controllers.ImageController{})

	beego.Router("/pornhub/*.html", &controllers.PornhubController{})
	beego.Router("/redtube/*.html", &controllers.RedtubeController{})
	beego.Router("/tube8/*.html", &controllers.Tube8Controller{})
	beego.Router("/youporn/*.html", &controllers.YoupornController{})
	beego.Router("/xtube/*.html", &controllers.XtubeController{})
	beego.Router("/spankwire/*.html", &controllers.SpankwireController{})
	beego.Router("/keezmovies/*.html", &controllers.KeezmoviesController{})
	beego.Router("/extremetube/*.html", &controllers.ExtremetubeController{})
}
