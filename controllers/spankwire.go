package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const SpankwireApiURL = "http://www.spankwire.com/api/HubTrafficApiCall"
const SpankwireApiTimeout = 2

type SpankwireVideo map[string]interface{}
type SpankwireEmbedCode map[string]interface{}

type SpankwireController struct {
	beego.Controller
}

func (c *SpankwireController) Get() {
	aux := strings.Replace(c.Ctx.Request.URL.Path, ".html", "", -1)
	str := strings.Split(aux, "/")
	videoID := str[2]

	redirect := "https://www.spankwire.com/title/video" + videoID + "?utm_source=just-tit.com&utm_medium=embed&utm_campaign=hubtraffic_dsmatilla"

	BaseDomain := "https://just-tit.com"

	type TemplateData = map[string]interface{}

	c.Data["ID"] = videoID
	c.Data["Domain"] = BaseDomain

	videocode := SpankwireGetVideoByID(videoID)
	video := videocode["video"].(map[string]interface{})
	embedcode := SpankwireGetVideoEmbedCode(videoID)
	embed := embedcode["embed"].(map[string]interface{})
	str2, _ := base64.StdEncoding.DecodeString(fmt.Sprintf("%s", embed["code"]))
	c.Data["Embed"] = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(string(str2))))
	c.Data["PageTitle"] = fmt.Sprintf("%s", video["title"])
	c.Data["PageMetaDesc"] = fmt.Sprintf("%s", video["title"])
	c.Data["Thumb"] = fmt.Sprintf("%s", video["thumb"])
	c.Data["Url"] = fmt.Sprintf(BaseDomain+"/spankwire/%s.html", videoID)
	c.Data["Width"] = "650"
	c.Data["Height"] = "550"
	c.Data["SpankwireVideo"] = video

	if c.Data["PageTitle"] == "" {
		c.Redirect(redirect, 307)
	}

	if c.GetString("tp") == "true" {
		c.TplName = "video/player.html"
	} else {
		c.Data["Result"] = doSearch(fmt.Sprintf("%s", fmt.Sprintf("%s", video["title"])))
		c.TplName = "index.html"
	}
}

func SpankwireGetVideoByID(ID string) SpankwireVideo {
	timeout := time.Duration(SpankwireApiTimeout * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(fmt.Sprintf(SpankwireApiURL+"?data=getVideoById&output=json&video_id=%s", ID))
	if err != nil {
		return SpankwireVideo{}
		log.Println("[SPANKWIRE][GETVIDEOBYID]",err)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	var result SpankwireVideo
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Println("[SPANKWIRE][GETVIDEOBYID]",err)
	}
	return result

}

func SpankwireGetVideoEmbedCode(ID string) SpankwireEmbedCode {
	timeout := time.Duration(SpankwireApiTimeout * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(fmt.Sprintf(SpankwireApiURL+"?data=getVideoEmbedCode&output=json&video_id=%s", ID))
	if err != nil {
		return SpankwireEmbedCode{}
		log.Println("[SPANKWIRE][GETVIDEOEMBEDCODE]",err)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	var result SpankwireEmbedCode
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Println("[SPANKWIRE][GETVIDEOEMBEDCODE]",err)
	}
	return result
}
