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

const redtubeAPIURL = "https://api.redtube.com/"
const redtubeAPITimeout = 3
const redtubeCacheDuration = time.Minute * 5

// RedtubeSearchResult type for redtube search result
type RedtubeSearchResult map[string]interface{}

// RedtubeSingleVideo type for redtube video result
type RedtubeSingleVideo map[string]interface{}

// RedtubeEmbedCode type for redtube embed code
type RedtubeEmbedCode map[string]interface{}

// RedtubeController Beego Controler
type RedtubeController struct {
	beego.Controller
}

// Get Redtube Video controller
func (c *RedtubeController) Get() {
	// Get videoID from URL
	aux := strings.Replace(c.Ctx.Request.URL.Path, ".html", "", -1)
	str := strings.Split(aux, "/")
	videoID := str[2]

	// Build redirect URL in case the API fails
	redirect := "https://www.redtube.com/" + videoID + "?utm_source=just-tit.com&utm_medium=embed&utm_campaign=hubtraffic_dsmatilla"

	// Get base domain from URL
	BaseDomain := "https://" + c.Controller.Ctx.Input.Domain()
	if c.Controller.Ctx.Input.Port() != 80 && c.Controller.Ctx.Input.Port() != 443 {
		BaseDomain += fmt.Sprintf("%s%d", ":", c.Controller.Ctx.Input.Port())
	}

	// Call the API and 307 redirect to fallback URL if something is not right
	data := redtubeGetVideoByID(videoID)
	_, ok := data["video"]
	if !ok {
		c.Redirect(redirect, 307)
		return
	}

	// Get Embed Code from API
	embedcode := redtubeGetVideoEmbedCode(videoID)
	if embedcode["embed"] == nil {
		c.Redirect(redirect, 307)
		return
	}
	temp := embedcode["embed"].(map[string]interface{})
	embed, _ := base64.StdEncoding.DecodeString(fmt.Sprintf("%s", temp["code"]))
	result := []JTVideo{}
	// Construct video object
	v := data["video"].(map[string]interface{})
	video := JTVideo{}
	video.ID = videoID
	video.Provider = "redtube"
	video.Domain = template.URL(BaseDomain)
	video.Title = fmt.Sprintf("%s", v["title"])
	video.Description = fmt.Sprintf("%s", v["title"])
	video.Thumb = fmt.Sprintf("%s", v["thumb"])
	video.Embed = template.HTML(fmt.Sprintf("<object><embed src=\"%+v\" /></object>", html.UnescapeString(string(embed))))
	video.URL = template.URL(fmt.Sprintf(BaseDomain+"/redtube/%s.html", videoID))
	video.Width = fmt.Sprintf("%s", v["width"])
	video.Height = fmt.Sprintf("%s", v["height"])
	video.Duration = fmt.Sprintf("%s", v["duration"])
	video.Views = fmt.Sprintf("%s", v["views"])
	video.Rating = fmt.Sprintf("%s", v["rating"])
	video.Ratings = fmt.Sprintf("%s", v["ratings"])
	video.Segment = fmt.Sprintf("%s", v["segment"])
	video.PublishDate = fmt.Sprintf("%s", v["publish_date"])
	video.Type = "single"
	for _, tags := range v["tags"].([]interface{}) {
		video.Tags = append(video.Tags, fmt.Sprintf("%s", tags))
	}

	for _, thumbs := range v["thumbs"].([]interface{}) {
		video.Thumbs = append(video.Thumbs, fmt.Sprintf("%s", thumbs.(map[string]interface{})["src"]))
	}

	video.ExternalID = fmt.Sprintf("%s", v["video_id"])
	video.ExternalURL = fmt.Sprintf("%s", v["url"])

	result = append(result, video)

	// Send object to template
	c.Data["PageTitle"] = video.Title
	c.Data["PageMetaDesc"] = video.Title
	c.Data["Result"] = result

	c.Data["SearchResult"] = doSearch(video.Title)

	if c.GetString("tp") == "true" {
		c.TplName = "player.tpl"
	} else {
		c.Layout = "index.tpl"
		c.TplName = "singlevideo.tpl"
	}
}

func redtubeGetVideoByID(ID string) RedtubeSingleVideo {
	Cached := JTCache.Get("redtube-video-"+ID)
	var result RedtubeSingleVideo
	if Cached == nil {
		timeout := time.Duration(redtubeAPITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(redtubeAPIURL+"?data=redtube.Videos.getVideoById&video_id=%s&output=json", ID))
		if err != nil {
			log.Println("[REDTUBE][GETVIDEOBYID]",err)
			return RedtubeSingleVideo{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[REDTUBE][GETVIDEOBYID]",err)
			return RedtubeSingleVideo{}
		}
		JTCache.Put("redtube-video-"+ID, b, redtubeCacheDuration)
	} else {
		json.Unmarshal(Cached.([]uint8), &result)
	}

	return result
}

func redtubeGetVideoEmbedCode(ID string) RedtubeEmbedCode {
	Cached := JTCache.Get("redtube-embed-"+ID)
	if Cached == nil {
		timeout := time.Duration(redtubeAPITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(redtubeAPIURL+"?data=redtube.Videos.getVideoEmbedCode&video_id=%s&output=json", ID))
		if err != nil {
			log.Println("[REDTUBE][GETVIDEOEMBEDCODE]",err)
			return RedtubeEmbedCode{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result RedtubeEmbedCode
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[REDTUBE][GETVIDEOEMBEDCODE]",err)
			return RedtubeEmbedCode{}
		}
		JTCache.Put("redtube-embed-"+ID, b, redtubeCacheDuration)
		return result
	}
	var result RedtubeEmbedCode
	json.Unmarshal(Cached.([]uint8), &result)
	return result
}

// RedtubeSearch Calls redtube search function and process result to get array of videos
func RedtubeSearch(search string) []JTVideo {
	videos := redtubeSearchVideos(search)
	result := []JTVideo{}
	if videos["videos"] != nil {
		for _, data := range videos["videos"].([]interface{}) {
			// Construct video object
			v := data.(map[string]interface{})["video"]
			video := JTVideo{}
			video.ID = fmt.Sprintf("%s", v.(map[string]interface{})["video_id"])
			video.Provider = "redtube"
			video.Title = fmt.Sprintf("%s", v.(map[string]interface{})["title"])
			video.Description = fmt.Sprintf("%s", v.(map[string]interface{})["title"])
			video.Thumb = fmt.Sprintf("%s", v.(map[string]interface{})["thumb"])
			video.Width = fmt.Sprintf("%s", v.(map[string]interface{})["width"])
			video.Height = fmt.Sprintf("%s", v.(map[string]interface{})["height"])
			video.Duration = fmt.Sprintf("%s", v.(map[string]interface{})["duration"])
			video.Views = fmt.Sprintf("%.0f", v.(map[string]interface{})["views"])
			video.Rating = fmt.Sprintf("%.0f", v.(map[string]interface{})["rating"])
			video.Ratings = fmt.Sprintf("%.0f", v.(map[string]interface{})["ratings"])
			video.Segment = fmt.Sprintf("%s", v.(map[string]interface{})["segment"])
			video.PublishDate = fmt.Sprintf("%s", v.(map[string]interface{})["publish_date"])
			video.ExternalID = fmt.Sprintf("%s", v.(map[string]interface{})["video_id"])
			video.ExternalURL = fmt.Sprintf("%s", v.(map[string]interface{})["url"])
			video.Type = "search"
			tags := v.(map[string]interface{})["tags"]
			for _, tag := range tags.([]interface{}) {
				video.Tags = append(video.Tags, fmt.Sprintf("%s", tag.(map[string]interface{})["tag_name"]))
			}

			thumbs := v.(map[string]interface{})["thumbs"]
			for _, thumb := range thumbs.([]interface{}) {
				video.Thumbs = append(video.Thumbs, fmt.Sprintf("%s", thumb.(map[string]interface{})["src"]))
			}

			result = append(result, video)
		}
	}
	return result
}

func redtubeSearchVideos(search string) RedtubeSearchResult {
	Cached := JTCache.Get("redtube-search-"+search)
	if Cached == nil {
		timeout := time.Duration(redtubeAPITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(redtubeAPIURL+"?data=redtube.Videos.searchVideos&output=json&search=%s&thumbsize=small", url.QueryEscape(search)))
		if err != nil {
			log.Println("[REDTUBE][SEARCHVIDEOS]",err)
			return RedtubeSearchResult{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result RedtubeSearchResult
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[REDTUBE][SEARCHVIDEOS]",err)
			return RedtubeSearchResult{}
		}
		JTCache.Put("redtube-search-"+search, b, redtubeCacheDuration)
		return result
	}
	var result RedtubeSearchResult
	json.Unmarshal(Cached.([]uint8), &result)
	return result
}