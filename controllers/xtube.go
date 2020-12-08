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

const xtubeAPIURL = "http://www.xtube.com/webmaster/api.php"
const xtubeAPITimeout = 3
const xtubeCacheDuration = time.Minute * 5

// XtubeSingleVideo type for xtube api video
type XtubeSingleVideo map[string]interface{}

// XtubeSearchResult type for xtube api search result
type XtubeSearchResult []interface{}

// XtubeController Beego Controller
type XtubeController struct {
	beego.Controller
}

// Get Xtube Video Controller
func (c *XtubeController) Get() {
	aux := strings.Replace(c.Ctx.Request.URL.Path, ".html", "", -1)
	str := strings.Split(aux, "/")
	videoID := str[2]

	// Build redirect URL in case the API fails
	redirect := "https://www.xtube.com/video-watch/watchin-xtube-" + videoID + "?t=0&utm_source=just-tit.com&utm_medium=embed&utm_campaign=hubtraffic_dsmatilla"

	// Get base domain from URL
	BaseDomain := "https://" + c.Controller.Ctx.Input.Domain()
	if c.Controller.Ctx.Input.Port() != 80 && c.Controller.Ctx.Input.Port() != 443 {
		BaseDomain += fmt.Sprintf("%s%d", ":", c.Controller.Ctx.Input.Port())
	}

	// Call the API and 307 redirect to fallback URL if something is not right
	data := xtubeGetVideoByID(videoID)
	_, ok := data["video_id"]
	if !ok {
		c.Redirect(redirect, 307)
		return
	}

	result := []JTVideo{}
	// Construct video object
	v := data
	video := JTVideo{}
	video.ID = videoID
	video.Provider = "xtube"
	video.Domain = template.URL(BaseDomain)
	video.Title = fmt.Sprintf("%s", v["title"])
	video.Description = fmt.Sprintf("%s", v["title"])
	video.Thumb = fmt.Sprintf("%s", v["thumb"])
	video.Embed = template.HTML(fmt.Sprintf("<object><embed src=\"%+v\" /></object>", html.UnescapeString(v["embedCode"].(string))))
	video.URL = template.URL(fmt.Sprintf(BaseDomain+"/xtube/%s.html", videoID))
	video.Width = fmt.Sprintf("%s", v["width"])
	video.Height = fmt.Sprintf("%s", v["height"])
	video.Duration = fmt.Sprintf("%s", v["duration"])
	video.Views = fmt.Sprintf("%s", v["views"])
	video.Rating = fmt.Sprintf("%s", v["rating"])
	video.Ratings = fmt.Sprintf("%s", v["ratings"])
	video.Segment = fmt.Sprintf("%s", v["segment"])
	video.PublishDate = fmt.Sprintf("%s", v["publish_date"])
	video.Type = "single"
	for _, tags := range v["tags"].(map[string]interface{}) {
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

func xtubeGetVideoByID(ID string) XtubeSingleVideo {
	Cached := JTCache.Get("xtube-video-" + ID)
	var result XtubeSingleVideo
	if Cached == nil {
		timeout := time.Duration(xtubeAPITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(xtubeAPIURL+"?action=getVideoById&video_id=%s", ID))
		if err != nil {
			log.Println("[XTUBE][GETVIDEOBYID]", err)
			return XtubeSingleVideo{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[XTUBE][GETVIDEOBYID]", err)
		}
		JTCache.Put("xtube-video-"+ID, b, xtubeCacheDuration)
	} else {
		json.Unmarshal(Cached.([]uint8), &result)
	}
	return result
}

// XtubeSearch Calls xtube search function and process result to get array of videos
func XtubeSearch(search string) []JTVideo {
	videos := xtubeSearchVideos(search)
	result := []JTVideo{}
	if videos != nil {
		for _, data := range videos {
			// Construct video object
			v := data
			video := JTVideo{}
			video.ID = fmt.Sprintf("%s", v.(map[string]interface{})["video_id"])
			video.Provider = "xtube"
			video.Title = fmt.Sprintf("%s", v.(map[string]interface{})["title"])
			video.Description = fmt.Sprintf("%s", v.(map[string]interface{})["title"])
			video.Thumb = fmt.Sprintf("%s", v.(map[string]interface{})["thumb"])
			video.Width = fmt.Sprintf("%s", v.(map[string]interface{})["width"])
			video.Height = fmt.Sprintf("%s", v.(map[string]interface{})["height"])
			video.Duration = fmt.Sprintf("%s", v.(map[string]interface{})["duration"])
			video.Views = fmt.Sprintf("%.0f", v.(map[string]interface{})["views"])
			video.Rating = fmt.Sprintf("%s", v.(map[string]interface{})["rating"])
			video.Ratings = fmt.Sprintf("%.0f", v.(map[string]interface{})["ratings"])
			video.Segment = fmt.Sprintf("%s", v.(map[string]interface{})["segment"])
			video.PublishDate = fmt.Sprintf("%s", v.(map[string]interface{})["publish_date"])
			video.ExternalID = fmt.Sprintf("%s", v.(map[string]interface{})["video_id"])
			video.ExternalURL = fmt.Sprintf("%s", v.(map[string]interface{})["url"])
			video.Type = "search"
			tags := v.(map[string]interface{})["tags"]
			for _, tag := range tags.(map[string]interface{}) {
				video.Tags = append(video.Tags, fmt.Sprintf("%s", tag))
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

func xtubeSearchVideos(search string) XtubeSearchResult {
	Cached := JTCache.Get("xtube-search-" + search)
	if Cached == nil {
		timeout := time.Duration(xtubeAPITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(xtubeAPIURL+"?action=getVideosBySearchParams&search=%s&period=lastweek&ordering=latest&count=10", url.QueryEscape(search)))
		if err != nil {
			log.Println("[XTUBE][SEARCHVIDEOS]", err)
			return XtubeSearchResult{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result XtubeSearchResult
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[XTUBE][SEARCHVIDEOS]", err)
			return XtubeSearchResult{}
		}
		JTCache.Put("xtube-search-"+search, b, xtubeCacheDuration)
		return result
	}
	var result XtubeSearchResult
	json.Unmarshal(Cached.([]uint8), &result)
	return result
}
