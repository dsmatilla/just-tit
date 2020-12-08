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

const tube8APIURL = "http://api.tube8.com/api.php"
const tube8APITimeout = 3
const tube8CacheDuration = time.Minute * 5

// Tube8SearchResult type for tube8 search result
type Tube8SearchResult map[string]interface{}

// Tube8SingleVideo type for tube8 search result
type Tube8SingleVideo map[string]interface{}

// Tube8Controller Beego Controler
type Tube8Controller struct {
	beego.Controller
}

// Get Tube8 Video controller
func (c *Tube8Controller) Get() {
	aux := strings.Replace(c.Ctx.Request.URL.Path, ".html", "", -1)
	str := strings.Split(aux, "/")
	videoID := str[2]

    // Build redirect URL in case the API fails
	redirect := "https://www.tube8.com/video/title/" + videoID + "/?utm_source=just-tit.com&utm_medium=embed&utm_campaign=hubtraffic_dsmatilla"

	// Get base domain from URL
	BaseDomain := "https://" + c.Controller.Ctx.Input.Domain()
	if c.Controller.Ctx.Input.Port() != 80 && c.Controller.Ctx.Input.Port() != 443 {
		BaseDomain += fmt.Sprintf("%s%d", ":", c.Controller.Ctx.Input.Port())
	}

	// Call the API and 307 redirect to fallback URL if something is not right
	data := tube8GetVideoByID(videoID)
	_, ok := data["video"]
	if !ok {
		c.Redirect(redirect, 307)
		return
	}

	// Get Embed Code from API
	embed := tube8GetVideoEmbedCode(videoID)
	if embed == "" {
		c.Redirect(redirect, 307)
		return
	}
	embed = strings.Replace(embed, "[\"", "", -1)
	embed = strings.Replace(embed, "\"]", "", -1)
	embed = strings.Replace(embed, "\\\"", "\"", -1)
	embed = strings.Replace(embed, "\\/", "/", -1)

	result := []JTVideo{}
	// Construct video object
	v := data["video"].(map[string]interface{})
	video := JTVideo{}
	video.ID = videoID
	video.Provider = "tube8"
	video.Domain = template.URL(BaseDomain)
	video.Title = fmt.Sprintf("%s", data["title"])
	video.Description = fmt.Sprintf("%s", data["title"])
	video.Thumb = fmt.Sprintf("%s", v["default_thumb"])
	video.Embed = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(embed)))
	video.URL = template.URL(fmt.Sprintf(BaseDomain+"/tube8/%s.html", videoID))
	video.Width = fmt.Sprintf("%s", v["width"])
	video.Height = fmt.Sprintf("%s", v["height"])
	video.Duration = fmt.Sprintf("%.0f", v["duration"])
	video.Views = fmt.Sprintf("%s", v["views"])
	video.Rating = fmt.Sprintf("%s", v["rating"])
	video.Ratings = fmt.Sprintf("%s", v["ratings"])
	video.Segment = fmt.Sprintf("%s", v["segment"])
	video.PublishDate = fmt.Sprintf("%s", v["publish_date"])
	video.Type = "single"
	for _, tags := range data["tags"].([]interface{}) {
		video.Tags = append(video.Tags, fmt.Sprintf("%s", tags))
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

func tube8GetVideoByID(ID string) Tube8SingleVideo {
	Cached := JTCache.Get("tube8-video-"+ID)
	var result Tube8SingleVideo
	if Cached == nil {
		timeout := time.Duration(tube8APITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(tube8APIURL+"?action=getvideobyid&video_id=%s&output=json&thumbsize=all", ID))
		if err != nil {
			log.Println("[TUBE8][GETVIDEOBYID]",err)
			return Tube8SingleVideo{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[TUBE8][GETVIDEOBYID]",err)
			return Tube8SingleVideo{}
		}
		JTCache.Put("tube8-video-"+ID, b, tube8CacheDuration)
	} else {
		json.Unmarshal(Cached.([]uint8), &result)
	}

	return result
}

func tube8GetVideoEmbedCode(ID string) string {
	Cached := JTCache.Get("tube8-embed-"+ID)
	if Cached == nil {
		timeout := time.Duration(tube8APITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(tube8APIURL+"?action=getvideoembedcode&output=json&video_id=%s", ID))
		if err != nil {
			log.Println("[TUBE8][GETVIDEOEMBEDCODE]",err)
			return ""
		}
		b, _ := ioutil.ReadAll(resp.Body)
		result := fmt.Sprintf("%s", b)
		JTCache.Put("tube8-embed-"+ID, b, tube8CacheDuration)
		return result
	}
	var result string
	json.Unmarshal(Cached.([]uint8), &result)
	return result
}

// Tube8Search Calls tube8 search function and process result to get array of videos
func Tube8Search(search string) []JTVideo {
	videos := tube8SearchVideos(search)

	result := []JTVideo{}
	if videos["videos"] != nil {		
		for _, data := range videos["videos"].([]interface{}) {
			// Construct video object
			v := data.(map[string]interface{})["video"]
			video := JTVideo{}
			video.ID = fmt.Sprintf("%.0f", v.(map[string]interface{})["video_id"])
			video.Provider = "tube8"
			video.Title = fmt.Sprintf("%s", data.(map[string]interface{})["title"])
			video.Description = fmt.Sprintf("%s", data.(map[string]interface{})["title"])
			video.Thumb = fmt.Sprintf("%s", v.(map[string]interface{})["default_thumb"])
			video.Width = fmt.Sprintf("%s", v.(map[string]interface{})["width"])
			video.Height = fmt.Sprintf("%s", v.(map[string]interface{})["height"])
			video.Duration = fmt.Sprintf("%.0f", v.(map[string]interface{})["duration"])
			video.Views = fmt.Sprintf("%.0f", v.(map[string]interface{})["views"])
			video.Rating = fmt.Sprintf("%s", v.(map[string]interface{})["rating"])
			video.Ratings = fmt.Sprintf("%.0f", v.(map[string]interface{})["ratings"])
			video.Segment = fmt.Sprintf("%s", v.(map[string]interface{})["segment"])
			video.PublishDate = fmt.Sprintf("%s", v.(map[string]interface{})["publish_date"])
			video.ExternalID = fmt.Sprintf("%s", v.(map[string]interface{})["video_id"])
			video.ExternalURL = fmt.Sprintf("%s", v.(map[string]interface{})["url"])
			video.Type = "search"
			tags := data.(map[string]interface{})["tags"]
			for _, tag := range tags.([]interface{}) {
				video.Tags = append(video.Tags, fmt.Sprintf("%s", tag))
			}

			thumbs := data.(map[string]interface{})["thumbs"]
			aux := thumbs.(map[string]interface{})["small"]
			for _, thumb := range aux.([]interface{}) {
				video.Thumbs = append(video.Thumbs, fmt.Sprintf("%s", thumb))
			}


			result = append(result, video)
		}
	}
	return result
}

func tube8SearchVideos(search string) Tube8SearchResult {
	Cached := JTCache.Get("tube8-search-"+search)
	if Cached == nil {
		timeout := time.Duration(tube8APITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(tube8APIURL+"?action=searchVideos&output=json&search=%s&thumbsize=all", url.QueryEscape(search)))
		if err != nil {
			log.Println("[TUBE8][SEARCHVIDEOS]",err)
			return Tube8SearchResult{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result Tube8SearchResult
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[TUBE8][SEARCHVIDEOS]",err)
			return Tube8SearchResult{}
		}
		JTCache.Put("tube8-search-"+search, b, tube8CacheDuration)
		return result
	}
	var result Tube8SearchResult
	json.Unmarshal(Cached.([]uint8), &result)
	return result
}