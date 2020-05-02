package controllers

import (
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"strings"
)

type ImageController struct {
	beego.Controller
}

func (c *ImageController) Get() {
	image := strings.Replace(c.Ctx.Request.URL.Path, "/images/", "", -1)
	aux := strings.Split(image, ".")
	str, _ := base64.StdEncoding.DecodeString(aux[0])
	response, _ := http.Get(fmt.Sprintf("%s", str))
	raw, _ := ioutil.ReadAll(response.Body)
	c.Ctx.Output.Header("Cache-Control", "max-age=31536000")
	c.Ctx.Output.Body(raw)
}
