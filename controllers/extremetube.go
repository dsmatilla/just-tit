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

const ExtremetubeApiURL = "https://www.extremetube.com/api/HubTrafficApiCall?"
const ExtremetubeApiTimeout = 5

type ExtremetubeEmbedCode map[string]interface{}
type ExtremetubeSingleVideo map[string]interface{}

type ExtremetubeController struct {
	beego.Controller
}

func (c *ExtremetubeController) Get() {
	aux := strings.Replace(c.Ctx.Request.URL.Path, ".html", "", -1)
	str := strings.Split(aux, "/")
	videoID := str[2]

	redirect := "https://www.extremetube.com/video/title-" + videoID + "?utm_source=just-tit.com&utm_medium=embed&utm_campaign=hubtraffic_dsmatilla"

	BaseDomain := "https://just-tit.com"

	type TemplateData = map[string]interface{}

	c.Data["ID"] = videoID
	c.Data["Domain"] = BaseDomain

	videocode := ExtremetubeGetVideoByID(videoID)
	_, ok := videocode["video"]
	if !ok { c.Redirect(redirect, 307) }
	video := videocode["video"].(map[string]interface{})
	embedcode := ExtremetubeGetVideoEmbedCode(videoID)
	embed := embedcode["embed"].(map[string]interface{})

	str2, _ := base64.StdEncoding.DecodeString(fmt.Sprintf("%s", embed["code"]))
	c.Data["Embed"] = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(string(str2))))
	c.Data["PageTitle"] = fmt.Sprintf("%s", video["title"])
	c.Data["PageMetaDesc"] = fmt.Sprintf("%s", video["title"])
	c.Data["Thumb"] = fmt.Sprintf("%s", video["default_thumb"])
	c.Data["Url"] = fmt.Sprintf(BaseDomain+"/extremetube/%s.html", videoID)
	c.Data["Width"] = "650"
	c.Data["Height"] = "550"
	c.Data["ExtremetubeVideo"] = video

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

func ExtremetubeGetVideoByID(ID string) ExtremetubeSingleVideo {
	timeout := time.Duration(ExtremetubeApiTimeout * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, _ := client.Get(fmt.Sprintf(ExtremetubeApiURL+"data=getVideoById&output=json&video_id=%s", ID))
	b, _ := ioutil.ReadAll(resp.Body)
	var result ExtremetubeSingleVideo
	err := json.Unmarshal(b, &result)
	if err != nil {
		log.Println("[EXTREMETUBE][GETVIDEOBYID]",err)
	}
	return result

}

func ExtremetubeGetVideoEmbedCode(ID string) ExtremetubeEmbedCode {
	timeout := time.Duration(ExtremetubeApiTimeout * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, _ := client.Get(fmt.Sprintf(ExtremetubeApiURL+"data=getVideoEmbedCode&video_id=%s", ID))
	b, _ := ioutil.ReadAll(resp.Body)
	var result ExtremetubeEmbedCode
	err := json.Unmarshal(b, &result)
	if err != nil {
		log.Println("[EXTREMETUBE][GETVIDEOEMBEDCODE]",err)
	}
	return result
}
