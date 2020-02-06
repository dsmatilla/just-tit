package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const Tube8ApiURL = "http://api.tube8.com/api.php"
const Tube8ApiTimeout = 2
const Tube8CacheDuration = time.Minute * 5

type Tube8SearchResult map[string]interface{}
type Tube8SingleVideo map[string]interface{}

type Tube8Controller struct {
	beego.Controller
}

func (c *Tube8Controller) Get() {
	aux := strings.Replace(c.Ctx.Request.URL.Path, ".html", "", -1)
	str := strings.Split(aux, "/")
	videoID := str[2]

	redirect := "https://www.tube8.com/video/title/" + videoID + "/?utm_source=just-tit.com&utm_medium=embed&utm_campaign=hubtraffic_dsmatilla"

	BaseDomain := "https://just-tit.com"

	type TemplateData = map[string]interface{}

	c.Data["ID"] = videoID
	c.Data["Domain"] = BaseDomain

	videocode := Tube8GetVideoByID(videoID)
	video := videocode["video"].(map[string]interface{})
	embed := Tube8GetVideoEmbedCode(videoID)
	embed = strings.Replace(embed, "[\"", "", -1)
	embed = strings.Replace(embed, "\"]", "", -1)
	embed = strings.Replace(embed, "\\\"", "\"", -1)
	embed = strings.Replace(embed, "\\/", "/", -1)
	c.Data["Embed"] = template.HTML(fmt.Sprintf("%+v", embed))
	c.Data["PageTitle"] = fmt.Sprintf("%s", videocode["title"])
	c.Data["PageMetaDesc"] = fmt.Sprintf("%s", videocode["title"])
	c.Data["Thumb"] = fmt.Sprintf("%s", video["default_thumb"])
	c.Data["Url"] = fmt.Sprintf(BaseDomain+"/tube8/%s.html", videoID)
	c.Data["Width"] = "628"
	c.Data["Height"] = "362"
	c.Data["Tube8Video"] = videocode

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

func Tube8SearchVideos(search string) Tube8SearchResult {
	Cached := JTCache.Get("tube8-search-"+search)
	if Cached == nil {
		timeout := time.Duration(Tube8ApiTimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(Tube8ApiURL+"?action=searchVideos&output=json&search=%s&thumbsize=all", url.QueryEscape(search)))
		if err != nil {
			return Tube8SearchResult{}
			log.Println("[TUBE8][SEARCHVIDEOS]",err)
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result Tube8SearchResult
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[TUBE8][SEARCHVIDEOS]",err)
		}
		JTCache.Put("tube8-search-"+search, result, Tube8CacheDuration)
		return result
	} else {
		return Cached.(Tube8SearchResult)
	}
}

func Tube8GetVideoByID(ID string) Tube8SingleVideo {
	Cached := JTCache.Get("tube8-video-"+ID)
	if Cached == nil {
		timeout := time.Duration(Tube8ApiTimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(Tube8ApiURL+"?action=getvideobyid&video_id=%s&output=json&thumbsize=all", ID))
		if err != nil {
			return Tube8SingleVideo{}
			log.Println("[TUBE8][GETVIDEOBYID]",err)
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result Tube8SingleVideo
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[TUBE8][GETVIDEOBYID]",err)
		}
		JTCache.Put("tube8-video-"+ID, result, Tube8CacheDuration)
		return result
	} else {
		return Cached.(Tube8SingleVideo)
	}
}

func Tube8GetVideoEmbedCode(ID string) string {
	Cached := JTCache.Get("tube8-embed-"+ID)
	if Cached == nil {
		timeout := time.Duration(Tube8ApiTimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(Tube8ApiURL+"?action=getvideoembedcode&output=json&video_id=%s", ID))
		if err != nil {
			log.Println("[TUBE8][GETVIDEOEMBEDCODE]",err)
		}
		b, _ := ioutil.ReadAll(resp.Body)
		result := fmt.Sprintf("%s", b)
		JTCache.Put("tube8-embed-"+ID, result, Tube8CacheDuration)
		return result
	} else {
		return Cached.(string)
	}
}

