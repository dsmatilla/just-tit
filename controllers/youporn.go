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

const YoupornApiURL = "http://www.youporn.com/api/webmasters/"
const YoupornApiTimeout = 2
const YoupornCacheDuration = time.Minute * 10

type YoupornSearchResult map[string]interface{}
type YoupornSingleVideo map[string]interface{}
type YoupornEmbedCode map[string]interface{}

type YoupornController struct {
	beego.Controller
}

func (c *YoupornController) Get() {
	aux := strings.Replace(c.Ctx.Request.URL.Path, ".html", "", -1)
	str := strings.Split(aux, "/")
	videoID := str[2]

	redirect := "https://www.youporn.com/watch/" + videoID + "/title/?utm_source=just-tit.com&utm_medium=embed&utm_campaign=hubtraffic_dsmatilla"

	BaseDomain := "https://"+c.Ctx.Request.Host

	type TemplateData = map[string]interface{}

	c.Data["ID"] = videoID
	c.Data["Domain"] = BaseDomain

	videocode := YoupornGetVideoByID(videoID)
	video := videocode["video"].(map[string]interface{})
	embed := YoupornGetVideoEmbedCode(videoID)
	c.Data["Embed"] = template.HTML(fmt.Sprintf("%+v", embed))
	c.Data["PageTitle"] = fmt.Sprintf("%s", video["title"])
	c.Data["PageMetaDesc"] = fmt.Sprintf("%s", video["title"])
	c.Data["Thumb"] = fmt.Sprintf("%s", video["default_thumb"])
	c.Data["Url"] = fmt.Sprintf(BaseDomain+"/youporn/%s.html", videoID)
	c.Data["Width"] = "628"
	c.Data["Height"] = "501"
	c.Data["YoupornVideo"] = videocode

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

func YoupornSearchVideos(search string) YoupornSearchResult {
	Cached := JTCache.Get("youporn-search-"+search)
	if Cached == nil {
		log.Println("SEARCH NOT CACHED")
		timeout := time.Duration(YoupornApiTimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(YoupornApiURL+"search?search=%s&thumbsize=all", url.QueryEscape(search)))
		if err != nil {
			return YoupornSearchResult{}
			log.Println("[YOUPORN][SEARCHVIDEOS]",err)
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result YoupornSearchResult
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[YOUPORN][SEARCHVIDEOS]",err)
		}
		JTCache.Put("youporn-search-"+search, result, YoupornCacheDuration)
		return result
	} else {
		log.Println("SEARCH CACHED")
		return Cached.(YoupornSearchResult)
	}
}

func YoupornGetVideoByID(ID string) YoupornSingleVideo {
	Cached := JTCache.Get("youporn-video-"+ID)
	if Cached == nil {
		log.Println("VIDEO NOT CACHED")
		timeout := time.Duration(YoupornApiTimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(YoupornApiURL+"video_by_id/?video_id=%s", ID))
		if err != nil {
			return YoupornSingleVideo{}
			log.Println("[YOUPORN][GETVIDEOBYID]",err)
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result YoupornSingleVideo
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[YOUPORN][GETVIDEOBYID]",err)
		}
		JTCache.Put("youporn-video-"+ID, result, YoupornCacheDuration)
		return result
	} else {
		log.Println("VIDEO CACHED")
		return Cached.(YoupornSingleVideo)
	}
}

func YoupornGetVideoEmbedCode(ID string) YoupornEmbedCode {
	Cached := JTCache.Get("youporn-embed-"+ID)
	if Cached == nil {
		log.Println("EMBED NOT CACHED")
		timeout := time.Duration(YoupornApiTimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(YoupornApiURL+"video_embed_code/?video_id=%s", ID))
		if err != nil {
			return YoupornEmbedCode{}
			log.Println("[YOUPORN][GETVIDEOEMBEDCODE]",err)
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result YoupornEmbedCode
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[YOUPORN][GETVIDEOEMBEDCODE]",err)
		}
		JTCache.Put("youporn-embed-"+ID, result, YoupornCacheDuration)
		return result
	} else {
		log.Println("EMBED CACHED")
		return Cached.(YoupornEmbedCode)
	}
}