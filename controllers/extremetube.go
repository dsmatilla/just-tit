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
	"strconv"
)

const extremetubeAPIURL = "https://www.extremetube.com/api/HubTrafficApiCall"
const extremetubeAPITimeout = 3
const extremetubeCacheDuration = time.Minute * 5

// ExtremetubeSearchResult type for extremetube api search result
type ExtremetubeSearchResult map[string]interface{}

// ExtremetubeEmbedCode type for extremetube api embed code
type ExtremetubeEmbedCode map[string]interface{}

// ExtremetubeSingleVideo type for extremetube api video
type ExtremetubeSingleVideo map[string]interface{}

// ExtremetubeController Beego Controller
type ExtremetubeController struct {
	beego.Controller
}

// Get Extremetube Video Controller 
func (c *ExtremetubeController) Get() {
	aux := strings.Replace(c.Ctx.Request.URL.Path, ".html", "", -1)
	str := strings.Split(aux, "/")
	videoID := str[2]

	redirect := "https://www.extremetube.com/video/title-" + videoID + "?utm_source=just-tit.com&utm_medium=embed&utm_campaign=hubtraffic_dsmatilla"

	// Get base domain from URL
	BaseDomain := "https://" + c.Controller.Ctx.Input.Domain()
	if c.Controller.Ctx.Input.Port() != 80 && c.Controller.Ctx.Input.Port() != 443 {
		BaseDomain += fmt.Sprintf("%s%d", ":", c.Controller.Ctx.Input.Port())
	}

	// Call the API and 307 redirect to fallback URL if something is not right
	data := extremetubeGetVideoByID(videoID)
	_, ok := data["video"]
	if !ok {
		c.Redirect(redirect, 307)
		return
	}

	// Get Embed Code from API
	embedcode := extremetubeGetVideoEmbedCode(videoID)
	if embedcode["embed"] == nil {
		c.Redirect(redirect, 307)
		return
	}
	embed := embedcode["embed"].(map[string]interface{})

	result := []JTVideo{}
	// Construct video object
	v := data["video"].(map[string]interface{})
	video := JTVideo{}
	video.ID = videoID
	video.Provider = "extremetube"
	video.Domain = template.URL(BaseDomain)
	video.Title = fmt.Sprintf("%s", v["title"])
	video.Description = fmt.Sprintf("%s", v["title"])
	video.Thumb = fmt.Sprintf("%s", v["thumb"])
	str2, _ := base64.StdEncoding.DecodeString(fmt.Sprintf("%s", embed["code"]))
	video.Embed = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(string(str2))))
	video.URL = template.URL(fmt.Sprintf(BaseDomain+"/extremetube/%s.html", videoID))
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
		video.Tags = append(video.Tags, fmt.Sprintf("%s", tags.(map[string]interface{})["tag_name"]))
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

	if c.GetString("tp") == "true" {
		c.TplName = "player.tpl"
	} else {
		c.Layout = "index.tpl"
		c.TplName = "singlevideo.tpl"
	}
}

func extremetubeGetVideoByID(ID string) ExtremetubeSingleVideo {
	Cached := JTCache.Get("extremetube-video-" + ID)
	var result ExtremetubeSingleVideo
	if Cached == nil {
		timeout := time.Duration(extremetubeAPITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, _ := client.Get(fmt.Sprintf(extremetubeAPIURL+"?data=getVideoById&output=json&video_id=%s", ID))
		b, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[EXTREMETUBE][GETVIDEOBYID]",err)
			return ExtremetubeSingleVideo{}
		}
		JTCache.Put("extremetube-video-"+ID, b, extremetubeCacheDuration)
	} else {
		json.Unmarshal(Cached.([]uint8), &result)
	}
	return result
}

func extremetubeGetVideoEmbedCode(ID string) ExtremetubeEmbedCode {
	Cached := JTCache.Get("pornhub-embed-" + ID)
	if Cached == nil {
		timeout := time.Duration(extremetubeAPITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, _ := client.Get(fmt.Sprintf(extremetubeAPIURL+"?data=getVideoEmbedCode&video_id=%s", ID))
		b, _ := ioutil.ReadAll(resp.Body)
		var result ExtremetubeEmbedCode
		err := json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[EXTREMETUBE][GETVIDEOEMBEDCODE]",err)
		}
		return result
	}
	var result ExtremetubeEmbedCode
	json.Unmarshal(Cached.([]uint8), &result)
	return result
}

// ExtremetubeSearch Calls extremetube search function and process result to get array of videos
func ExtremetubeSearch(search string) []JTVideo {
	videos := extremetubeSearchVideos(search)
	result := []JTVideo{}
	if videos["videos"] != nil {
		for _, data := range videos["videos"].([]interface{}) {
			// Construct video object
			v := data.(map[string]interface{})["video"]
			video := JTVideo{}
			video.ID = fmt.Sprintf("%s", v.(map[string]interface{})["id"])
			if _, err := strconv.Atoi(video.ID); err != nil {
				video.ID = fmt.Sprintf("%.0f", v.(map[string]interface{})["id"])
			}
			video.Provider = "extremetube"
			video.Title = fmt.Sprintf("%s", v.(map[string]interface{})["title"])
			video.Description = fmt.Sprintf("%s", v.(map[string]interface{})["title"])
			video.Thumb = fmt.Sprintf("%s", v.(map[string]interface{})["thumb"])
			video.Width = fmt.Sprintf("%s", v.(map[string]interface{})["width"])
			video.Height = fmt.Sprintf("%s", v.(map[string]interface{})["height"])
			video.Duration = fmt.Sprintf("%s", v.(map[string]interface{})["duration"])
			video.Views = fmt.Sprintf("%s", v.(map[string]interface{})["views"])
			video.Rating = fmt.Sprintf("%.0f", v.(map[string]interface{})["rating"])
			video.Ratings = fmt.Sprintf("%.0f", v.(map[string]interface{})["ratings"])
			video.Segment = fmt.Sprintf("%s", v.(map[string]interface{})["segment"])
			video.PublishDate = fmt.Sprintf("%s", v.(map[string]interface{})["publish_date"])
			video.ExternalID = fmt.Sprintf("%s", v.(map[string]interface{})["id"])
			if _, err := strconv.Atoi(video.ExternalID); err != nil {
				video.ExternalID = fmt.Sprintf("%.0f", v.(map[string]interface{})["id"])
			}
			video.ExternalURL = fmt.Sprintf("%s", v.(map[string]interface{})["url"])
			video.Type = "search"
			tags := v.(map[string]interface{})["tags"] 
			for _, tag := range tags.([]interface{}) {
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

func extremetubeSearchVideos(search string) ExtremetubeSearchResult {
	Cached := JTCache.Get("extremetube-search-" + search)
	if Cached == nil {
		timeout := time.Duration(extremetubeAPITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(extremetubeAPIURL+"?data=searchVideos&output=json&search=%s&thumbsize=small", url.QueryEscape(search)))
		if err != nil {
			log.Println("[EXTREMETUBE][SEARCHVIDEOS]", err)
			return ExtremetubeSearchResult{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result ExtremetubeSearchResult
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[EXTREMETUBE][SEARCHVIDEOS]", err)
			return ExtremetubeSearchResult{}
		}
		JTCache.Put("extremetube-search-"+search, b, extremetubeCacheDuration)
		return result
	}
	var result ExtremetubeSearchResult
	json.Unmarshal(Cached.([]uint8), &result)
	return result
}