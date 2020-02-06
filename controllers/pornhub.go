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
	"net/url"
	"strings"
	"time"
)

const PornhubApiURL = "http://www.pornhub.com/webmasters/"
const PornhubApiTimeout = 2
const PornhubCacheDuration = time.Minute * 5

type PornhubSearchResult map[string]interface{}
type PornhubSingleVideo map[string]interface{}
type PornhubEmbedCode map[string]interface{}

type PornhubController struct {
	beego.Controller
}

func (c *PornhubController) Get() {
	aux := strings.Replace(c.Ctx.Request.URL.Path, ".html", "", -1)
	str := strings.Split(aux, "/")
	videoID := str[2]

	redirect := "https://pornhub.com/view_video.php?viewkey=" + videoID + "&t=1&utm_source=just-tit.com&utm_medium=embed&utm_campaign=hubtraffic_dsmatilla"

	BaseDomain := "https://just-tit.com"

	type TemplateData = map[string]interface{}

	c.Data["ID"] = videoID
	c.Data["Domain"] = BaseDomain

	videocode := PornhubGetVideoByID(videoID)
	video := videocode["video"].(map[string]interface{})
	embedcode := PornhubGetVideoEmbedCode(videoID)
	embed := embedcode["embed"].(map[string]interface{})
	c.Data["Embed"] = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(embed["code"].(string))))
	c.Data["PageTitle"] = fmt.Sprintf("%s", video["title"])
	c.Data["PageMetaDesc"] = fmt.Sprintf("%s", video["title"])
	c.Data["Thumb"] = fmt.Sprintf("%s", video["thumb"])
	c.Data["Url"] = fmt.Sprintf(BaseDomain+"/pornhub/%s.html", videoID)
	c.Data["Width"] = "580"
	c.Data["Height"] = "360"
	c.Data["PornhubVideo"] = videocode

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

func PornhubSearchVideos(search string) PornhubSearchResult {
	Cached := JTCache.Get("pornhub-search-"+search)
	if Cached == nil {
		timeout := time.Duration(PornhubApiTimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(PornhubApiURL+"search?search=%s&thumbnail=all", url.QueryEscape(search)))
		if err != nil {
			log.Println("[PORNHUB][SEARCHVIDEOS]",err)
			return PornhubSearchResult{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result PornhubSearchResult
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[PORNHUB][SEARCHVIDEOS]",err)
		}
		JTCache.Put("pornhub-search-"+search, result, PornhubCacheDuration)
		return result
	} else {
		return Cached.(PornhubSearchResult)
	}
}

func PornhubGetVideoByID(ID string) PornhubSingleVideo {
	Cached := JTCache.Get("pornhub-video-"+ID)
	if Cached == nil {
		timeout := time.Duration(PornhubApiTimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(PornhubApiURL+"video_by_id?id=%s", ID))
		if err != nil {
			log.Println("[PORNHUB][GETVIDEOBYID]", err)
			return PornhubSingleVideo{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result PornhubSingleVideo
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[PORNHUB][GETVIDEOBYID]", err)
		}
		JTCache.Put("pornhub-video-"+ID, result, PornhubCacheDuration)
		return result
	} else {
		return Cached.(PornhubSingleVideo)
	}
}

func PornhubGetVideoEmbedCode(ID string) PornhubEmbedCode {
	Cached := JTCache.Get("pornhub-embed-"+ID)
	if Cached == nil {
		timeout := time.Duration(PornhubApiTimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(PornhubApiURL+"video_embed_code?id=%s", ID))
		if err != nil {
			log.Println("[PORNHUB][GETVIDEOEMBEDCODE]",err)
			return PornhubEmbedCode{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result PornhubEmbedCode
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[PORNHUB][GETVIDEOEMBEDCODE]",err)
		}
		JTCache.Put("pornhub-embed-"+ID, result, PornhubCacheDuration)
		return result
	} else {
		return Cached.(PornhubEmbedCode)
	}
}
