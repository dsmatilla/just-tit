package main

import (
	"encoding/base64"
	"github.com/astaxie/beego"
	_ "github.com/dsmatilla/just-tit/routers"
	"strings"
)

func ToImageProxy(url string) string {
	aux := strings.Split(url, ".")
	ext := aux[len(aux)-1]
	if ext == "jpg" || ext == "png" {
		return "/images/" + base64.StdEncoding.EncodeToString([]byte(url)) + "." + ext
	} else {
		return url
	}
}

func main() {
	beego.AddFuncMap("ToImageProxy",ToImageProxy)
	beego.SetStaticPath("/img", "static/img")
	beego.SetStaticPath("/robots.txt", "static/robots.txt")
	beego.SetStaticPath("/service-worker.js", "static/js/service-worker.js")
	beego.SetStaticPath("/manifest.json", "static/manifest.json")

	beego.Run()
}

