package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/dsmatilla/just-tit/controllers"
	"os"
)

func init() {
	// Initialize memory cache
	redisHost := os.Getenv("redisHost")
	redisName := os.Getenv("redisName")
	redisDBNum := os.Getenv("redisDBNum")
	redisPasswd := os.Getenv("redisPasswd")
	if redisHost != "" {
		controllers.JTCache, _ = cache.NewCache("memory", `{"key":"`+redisName+`","conn":"`+redisHost+`","dbNum":"`+redisDBNum+`","password":"`+redisPasswd+`"}`)
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
