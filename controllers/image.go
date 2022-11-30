package controllers

import (
	"encoding/base64"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"io/ioutil"
	"net/http"
	"strings"
)

// ImageController beego controller
type ImageController struct {
	beego.Controller
}

// Get ImageController
func (c *ImageController) Get() {
	image := strings.Replace(c.Ctx.Request.URL.Path, "/images/", "", -1)
	aux := strings.Split(image, ".")
	str, _ := base64.StdEncoding.DecodeString(aux[0])
	response, _ := http.Get(fmt.Sprint(str))
	raw, _ := ioutil.ReadAll(response.Body)
	ct := http.DetectContentType(raw)
	c.Ctx.Output.Header("Cache-Control", "max-age=31536000")
	c.Ctx.Output.Header("Content-Type", ct)
	c.Ctx.Output.Body(raw)
}
