package controllers

import (
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

const XtubeApiURL = "http://www.xtube.com/webmaster/api.php"
const XtubeApiTimeout = 5

type XtubeVideo map[string]interface{}

type XtubeController struct {
	beego.Controller
}

func (c *XtubeController) Get() {
	aux := strings.Replace(c.Ctx.Request.URL.Path, ".html", "", -1)
	str := strings.Split(aux, "/")
	videoID := str[2]

	redirect := "https://www.xtube.com/video-watch/watchin-xtube-" + videoID + "?t=0&utm_source=just-tit.com&utm_medium=embed&utm_campaign=hubtraffic_dsmatilla"

	BaseDomain := "https://just-tit.com"

	type TemplateData = map[string]interface{}

	c.Data["ID"] = videoID
	c.Data["Domain"] = BaseDomain

	video := XtubeGetVideoByID(videoID)
	c.Data["Embed"] = template.HTML(fmt.Sprintf("<object><embed src=\"%+v\" /></object>", html.UnescapeString(video["embedCode"].(string))))
	c.Data["PageTitle"] = fmt.Sprintf("%s", video["title"])
	c.Data["PageMetaDesc"] = fmt.Sprintf("%s", video["title"])
	c.Data["Thumb"] = fmt.Sprintf("%s", video["thumb"])
	c.Data["Url"] = fmt.Sprintf(BaseDomain+"/xtube/%s.html", videoID)
	c.Data["Width"] = "628"
	c.Data["Height"] = "501"
	c.Data["XtubeVideo"] = video

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

func XtubeGetVideoByID(ID string) XtubeVideo {
	timeout := time.Duration(XtubeApiTimeout * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(fmt.Sprintf(XtubeApiURL+"?action=getVideoById&video_id=%s", ID))
	if err != nil {
		return XtubeVideo{}
		log.Println("[XTUBE][GETVIDEOBYID]",err)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	var result XtubeVideo
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Println("[XTUBE][GETVIDEOBYID]",err)
	}
	return result
}
