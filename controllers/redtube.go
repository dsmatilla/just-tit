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
	"net/url"
	"strings"
	"time"
)

const RedtubeApiURL = "https://api.redtube.com/"
const RedtubeApiTimeout = 2
const RedtubeCacheDuration = time.Minute * 5

type RedtubeSearchResult map[string]interface{}
type RedtubeSingleVideo map[string]interface{}
type RedtubeEmbedCode map[string]interface{}

type RedtubeController struct {
	beego.Controller
}

func (c *RedtubeController) Get() {
	aux := strings.Replace(c.Ctx.Request.URL.Path, ".html", "", -1)
	str := strings.Split(aux, "/")
	videoID := str[2]

	redirect := "https://www.redtube.com/" + videoID + "?utm_source=just-tit.com&utm_medium=embed&utm_campaign=hubtraffic_dsmatilla"

	BaseDomain := "https://just-tit.com"

	type TemplateData = map[string]interface{}

	c.Data["ID"] = videoID
	c.Data["Domain"] = BaseDomain

	videocode := RedtubeGetVideoByID(videoID)
	video := videocode["video"].(map[string]interface{})
	embedcode := RedtubeGetVideoEmbedCode(videoID)
	embed := embedcode["embed"].(map[string]interface{})
	str2, _ := base64.StdEncoding.DecodeString(fmt.Sprintf("%s", embed["code"]))
	c.Data["Embed"] = template.HTML(fmt.Sprintf("<object><embed src=\"%+v\" /></object>", html.UnescapeString(string(str2))))
	c.Data["PageTitle"] = fmt.Sprintf("%s", video["title"])
	c.Data["PageMetaDesc"] = fmt.Sprintf("%s", video["title"])
	c.Data["Thumb"] = fmt.Sprintf("%s", video["thumb"])
	c.Data["Url"] = fmt.Sprintf(BaseDomain+"/redtube/%s.html", videoID)
	c.Data["Width"] = "320"
	c.Data["Height"] = "180"
	c.Data["RedtubeVideo"] = videocode

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

func RedtubeSearchVideos(search string) RedtubeSearchResult {
	Cached := JTCache.Get("redtube-search-"+search)
	if Cached == nil {
		timeout := time.Duration(RedtubeApiTimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(RedtubeApiURL+"?data=redtube.Videos.searchVideos&output=json&search=%s&thumbsize=all", url.QueryEscape(search)))
		if err != nil {
			return RedtubeSearchResult{}
			log.Println("[REDTUBE][SEARCHVIDEOS]",err)
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result RedtubeSearchResult
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[REDTUBE][SEARCHVIDEOS]",err)
		}
		JTCache.Put("redtube-search-"+search, result, RedtubeCacheDuration)
		return result
	} else {
		return Cached.(RedtubeSearchResult)
	}
}

func RedtubeGetVideoByID(ID string) RedtubeSingleVideo {
	Cached := JTCache.Get("redtube-video-"+ID)
	if Cached == nil {
		timeout := time.Duration(RedtubeApiTimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(RedtubeApiURL+"?data=redtube.Videos.getVideoById&video_id=%s&output=json", ID))
		if err != nil {
			return RedtubeSingleVideo{}
			log.Println("[REDTUBE][GETVIDEOBYID]",err)
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result RedtubeSingleVideo
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[REDTUBE][GETVIDEOBYID]",err)
		}
		JTCache.Put("redtube-video-"+ID, result, RedtubeCacheDuration)
		return result
	} else {
		return Cached.(RedtubeSingleVideo)
	}
}

func RedtubeGetVideoEmbedCode(ID string) RedtubeEmbedCode {
	Cached := JTCache.Get("redtube-embed-"+ID)
	if Cached == nil {
		timeout := time.Duration(RedtubeApiTimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(RedtubeApiURL+"?data=redtube.Videos.getVideoEmbedCode&video_id=%s&output=json", ID))
		if err != nil {
			return RedtubeEmbedCode{}
			log.Println("[REDTUBE][GETVIDEOEMBEDCODE]",err)
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result RedtubeEmbedCode
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[REDTUBE][GETVIDEOEMBEDCODE]",err)
		}
		JTCache.Put("redtube-embed-"+ID, result, RedtubeCacheDuration)
		return result
	} else {
		return Cached.(RedtubeEmbedCode)
	}
}
