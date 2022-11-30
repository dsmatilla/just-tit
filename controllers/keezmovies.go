package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

const keezmoviesAPIURL = "http://www.keezmovies.com/wapi/"
const keezmoviesAPITimeout = 10
const keezmoviesCacheDuration = time.Minute * 5

// KeezmoviesSearchResult type for keezmovies api search result
type KeezmoviesSearchResult map[string]interface{}

// KeezmoviesEmbedCode type for keezmovies api embed code
type KeezmoviesEmbedCode map[string]interface{}

// KeezmoviesSingleVideo type for keezmovies api single video result
type KeezmoviesSingleVideo map[string]interface{}

// KeezmoviesController Beego Controller
type KeezmoviesController struct {
	beego.Controller
}

// Get Keexmovies Video controller
func (c *KeezmoviesController) Get() {
	// Get videoID from URL
	aux := strings.Replace(c.Ctx.Request.URL.Path, ".html", "", -1)
	str := strings.Split(aux, "/")
	videoID := str[2]

	// Build redirect URL in case the API fails
	redirect := "https://www.keezmovies.com/video/title-" + videoID + "?utm_source=just-tit.com&utm_medium=embed&utm_campaign=hubtraffic_dsmatilla"

	// Get base domain from URL
	BaseDomain := "https://" + c.Controller.Ctx.Input.Domain()
	if c.Controller.Ctx.Input.Port() != 80 && c.Controller.Ctx.Input.Port() != 443 {
		BaseDomain += fmt.Sprintf("%s%d", ":", c.Controller.Ctx.Input.Port())
	}

	// Call the API and 307 redirect to fallback URL if something is not right
	data := keezmoviesGetVideoByID(videoID)
	_, ok := data["video"]
	if !ok {
		c.Redirect(redirect, 307)
		return
	}

	// Get Embed Code from API
	embedcode := keezmoviesGetVideoEmbedCode(videoID)
	vembedcode := embedcode["video"]
	if vembedcode == nil {
		c.Redirect(redirect, 307)
		return
	}
	embed := vembedcode.(map[string]interface{})["embed_code"]

	result := []JTVideo{}
	// Construct video object
	v := data["video"].(map[string]interface{})
	video := JTVideo{}
	video.ID = videoID
	video.Provider = "keezmovies"
	video.Domain = template.URL(BaseDomain)
	video.Title = fmt.Sprintf("%s", v["title"])
	video.Description = fmt.Sprintf("%s", v["title"])
	video.Thumb = fmt.Sprintf("%s", v["image_url"])
	str2, _ := base64.StdEncoding.DecodeString(fmt.Sprintf("%s", embed))
	video.Embed = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(string(str2))))
	video.URL = template.URL(fmt.Sprintf(BaseDomain+"/keezmovies/%s.html", videoID))
	video.Width = fmt.Sprintf("%s", v["width"])
	video.Height = fmt.Sprintf("%s", v["height"])
	video.Duration = fmt.Sprintf("%s", v["duration"])
	video.Views = fmt.Sprintf("%s", v["times_viewed"])
	video.Rating = fmt.Sprintf("%s", v["rating"])
	video.Ratings = fmt.Sprintf("%s", v["ratings"])
	video.Segment = fmt.Sprintf("%s", v["segment"])
	video.PublishDate = fmt.Sprintf("%s", v["publish_date"])
	video.Type = "single"
	if v["tags"] != nil {
		tags := v["tags"]
		for _, tag := range tags.(map[string]interface{}) {
			video.Tags = append(video.Tags, fmt.Sprintf("%s", tag))
		}
	}
	if v["categories"] != nil {
		categories := v["categories"]
		for _, category := range categories.(map[string]interface{}) {
			video.Categories = append(video.Categories, fmt.Sprintf("%s", category))
		}
	}
	aux2 := v["thumbs"]
	aux3 := aux2.(map[string]interface{})["large"]
	thumb := aux3.(map[string]interface{})["flipbook_path"]
	// values := aux2.(map[string]interface{})["flipbook_values"] // TODO Check if the API has a bug, unexpected value
	for i := 1; i < 5; i++ {
		video.Thumbs = append(video.Thumbs, strings.Replace(fmt.Sprintf("%s", thumb), "{index}", fmt.Sprint(i), 1))
	}
	if v["pornstars"] != nil {
		pornstars := v["pornstars"]
		for _, pornstar := range pornstars.(map[string]interface{}) {
			video.Pornstars = append(video.Pornstars, fmt.Sprintf("%s", pornstar))
		}
	}
	video.ExternalID = fmt.Sprintf("%s", v["id"])
	video.ExternalURL = fmt.Sprintf("%s", v["video_url"])

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

func keezmoviesGetVideoByID(ID string) KeezmoviesSingleVideo {
	Cached, _ := JTCache.Get(context.Background(), "keezmovies-video-" + ID)
	var result KeezmoviesSingleVideo
	if Cached == nil {
		timeout := time.Duration(keezmoviesAPITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(keezmoviesAPIURL+"getVideoById?output=json&video_id=%s", ID))
		if err != nil {
			log.Println("[KEEZMOVIES][GETVIDEOBYID]", err)
			return KeezmoviesSingleVideo{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[KEEZMOVIES][GETVIDEOBYID]", err)
			return KeezmoviesSingleVideo{}
		}
		JTCache.Put(context.Background(), "keezmovies-video-"+ID, b, keezmoviesCacheDuration)
	} else {
		json.Unmarshal(Cached.([]uint8), &result)
	}
	return result
}

func keezmoviesGetVideoEmbedCode(ID string) KeezmoviesEmbedCode {
	Cached, _ := JTCache.Get(context.Background(), "keezmovies-embed-" + ID)
	if Cached == nil {
		timeout := time.Duration(keezmoviesAPITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(keezmoviesAPIURL+"getVideoEmbedCode?output=json&video_id=%s", ID))
		if err != nil {
			log.Println("[KEEZMOVIES][GETVIDEOBYID]", err)
			return KeezmoviesEmbedCode{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result KeezmoviesEmbedCode
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[KEEZMOVIES][GETVIDEOEMBEDCODE]", err)
			return KeezmoviesEmbedCode{}
		}
		JTCache.Put(context.Background(), "keezmovies-embed-"+ID, b, keezmoviesCacheDuration)
		return result
	}
	var result KeezmoviesEmbedCode
	json.Unmarshal(Cached.([]uint8), &result)
	return result
}

// KeezmoviesSearch Calls keezmovies search function and process result to get array of videos
func KeezmoviesSearch(search string) []JTVideo {
	videos := keezmoviesSearchVideos(search)
	result := []JTVideo{}
	aux := videos["search"]
	if aux != nil {
		aux2 := aux.(map[string]interface{})["videos"]
		for _, data := range aux2.(map[string]interface{}) {
			// Construct video object
			v := data
			video := JTVideo{}
			video.ID = fmt.Sprintf("%.0f", v.(map[string]interface{})["id"])
			video.Provider = "keezmovies"
			video.Title = fmt.Sprintf("%s", v.(map[string]interface{})["title"])
			video.Description = fmt.Sprintf("%s", v.(map[string]interface{})["title"])
			video.Thumb = fmt.Sprintf("%s", v.(map[string]interface{})["image_url"])
			video.Width = fmt.Sprintf("%s", v.(map[string]interface{})["width"])
			video.Height = fmt.Sprintf("%s", v.(map[string]interface{})["height"])
			video.Duration = fmt.Sprintf("%s", v.(map[string]interface{})["duration"])
			video.Views = fmt.Sprintf("%.0f", v.(map[string]interface{})["times_viewed"])
			video.Rating = fmt.Sprintf("%.0f", v.(map[string]interface{})["rating"])
			video.Ratings = fmt.Sprintf("%s", v.(map[string]interface{})["ratings"])
			video.Segment = fmt.Sprintf("%s", v.(map[string]interface{})["segment"])
			video.PublishDate = fmt.Sprintf("%s", v.(map[string]interface{})["publish_date"])
			video.ExternalID = fmt.Sprintf("%s", v.(map[string]interface{})["video_id"])
			video.ExternalURL = fmt.Sprintf("%s", v.(map[string]interface{})["video_url"])
			video.Type = "search"
			if v.(map[string]interface{})["tags"] != nil {
				tags := v.(map[string]interface{})["tags"]
				for _, tag := range tags.(map[string]interface{}) {
					video.Tags = append(video.Tags, fmt.Sprintf("%s", tag))
				}
			}
			if v.(map[string]interface{})["categories"] != nil {
				categories := v.(map[string]interface{})["categories"]
				for _, category := range categories.(map[string]interface{}) {
					video.Categories = append(video.Categories, fmt.Sprintf("%s", category))
				}
			}
			aux := v.(map[string]interface{})["thumbs"]
			aux2 := aux.(map[string]interface{})["large"]
			thumb := aux2.(map[string]interface{})["flipbook_path"]
			// values := aux2.(map[string]interface{})["flipbook_values"] // TODO Check if the API has a bug, unexpected value
			for i := 1; i < 5; i++ {
				video.Thumbs = append(video.Thumbs, strings.Replace(fmt.Sprintf("%s", thumb), "{index}", fmt.Sprint(i), 1))
			}
			if v.(map[string]interface{})["pornstars"] != nil {
				pornstars := v.(map[string]interface{})["pornstars"]
				for _, pornstar := range pornstars.(map[string]interface{}) {
					video.Pornstars = append(video.Pornstars, fmt.Sprintf("%s", pornstar))
				}
			}
			result = append(result, video)
		}
	}
	return result
}

func keezmoviesSearchVideos(search string) KeezmoviesSearchResult {
	Cached, _ := JTCache.Get(context.Background(), "keezmovies-search-" + search)
	if Cached == nil {
		timeout := time.Duration(keezmoviesAPITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(keezmoviesAPIURL+"searchVideos?query=%s&thumbnail=all", url.QueryEscape(search)))

		if err != nil {
			log.Println("[KEEZMOVIES][SEARCHVIDEOS]", err)
			return KeezmoviesSearchResult{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result KeezmoviesSearchResult
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[KEEZMOVIES][SEARCHVIDEOS]", err)
			return KeezmoviesSearchResult{}
		}
		JTCache.Put(context.Background(), "keezmovies-search-"+search, b, keezmoviesCacheDuration)
		return result
	}
	var result KeezmoviesSearchResult
	json.Unmarshal(Cached.([]uint8), &result)
	return result
}
