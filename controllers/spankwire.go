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

const spankwireAPIURL = "http://www.spankwire.com/api/HubTrafficApiCall"
const spankwireAPITimeout = 3
const spankwireCacheDuration = time.Minute * 5

// SpankwireSearchResult type for spankwire api search result
type SpankwireSearchResult map[string]interface{}

// SpankwireSingleVideo type for spankwire api video result
type SpankwireSingleVideo map[string]interface{}

// SpankwireEmbedCode type for spankwire api embed code result
type SpankwireEmbedCode map[string]interface{}

// SpankwireController Beego Controller
type SpankwireController struct {
	beego.Controller
}

// Get Spankwire Video controller
func (c *SpankwireController) Get() {
	// Get videoID from URL
	aux := strings.Replace(c.Ctx.Request.URL.Path, ".html", "", -1)
	str := strings.Split(aux, "/")
	videoID := str[2]

	// Build redirect URL in case the API fails
	redirect := "https://www.spankwire.com/title/video" + videoID + "?utm_source=just-tit.com&utm_medium=embed&utm_campaign=hubtraffic_dsmatilla"
	// Get base domain from URL
	BaseDomain := "https://" + c.Controller.Ctx.Input.Domain()
	if c.Controller.Ctx.Input.Port() != 80 && c.Controller.Ctx.Input.Port() != 443 {
		BaseDomain += fmt.Sprintf("%s%d", ":", c.Controller.Ctx.Input.Port())
	}

	// Call the API and 307 redirect to fallback URL if something is not right
	data := spankwireGetVideoByID(videoID)
	_, ok := data["video"]
	if !ok {
		c.Redirect(redirect, 307)
		return
	}

	// Get Embed Code from API
	embedcode := spankwireGetVideoEmbedCode(videoID)
	vembedcode := embedcode["embed"]
	if vembedcode == nil {
		c.Redirect(redirect, 307)
		return
	}
	embed := vembedcode.(map[string]interface{})["code"]

	result := []JTVideo{}
	// Construct video object
	v := data["video"].(map[string]interface{})
	video := JTVideo{}
	video.ID = videoID
	video.Provider = "spankwire"
	video.Domain = template.URL(BaseDomain)
	video.Title = fmt.Sprintf("%s", v["title"])
	video.Description = fmt.Sprintf("%s", v["title"])
	video.Thumb = fmt.Sprintf("%s", v["default_thumb"])
	str2, _ := base64.StdEncoding.DecodeString(fmt.Sprintf("%s", embed))
	video.Embed = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(string(str2))))
	video.URL = template.URL(fmt.Sprintf(BaseDomain+"/spankwire/%s.html", videoID))
	video.Width = fmt.Sprintf("%s", v["width"])
	video.Height = fmt.Sprintf("%s", v["height"])
	video.Duration = fmt.Sprintf("%s", v["duration"])
	video.Views = fmt.Sprintf("%s", v["times_viewed"])
	video.Rating = fmt.Sprintf("%s", v["rating"])
	video.Ratings = fmt.Sprintf("%s", v["ratings"])
	video.Segment = fmt.Sprintf("%s", v["segment"])
	video.PublishDate = fmt.Sprintf("%s", v["publish_date"])
	video.Type = "single"
	for _, tag := range v["tags"].([]interface{}) {
		log.Print(tag)
		video.Tags = append(video.Tags, fmt.Sprintf("%s", tag))
	}
	for _, thumbs := range v["thumbs"].([]interface{}) {
		video.Thumbs = append(video.Thumbs, fmt.Sprintf("%s", thumbs.(map[string]interface{})["src"]))
	}

	video.ExternalID = fmt.Sprintf("%s", v["id"])
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

func spankwireGetVideoByID(ID string) SpankwireSingleVideo {
	Cached := JTCache.Get("spankwire-video-" + ID)
	var result SpankwireSingleVideo
	if Cached == nil {
		timeout := time.Duration(spankwireAPITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(spankwireAPIURL+"?data=getVideoById&output=json&video_id=%s", ID))
		if err != nil {
			log.Println("[SPANKWIRE][GETVIDEOBYID]",err)
			return SpankwireSingleVideo{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[SPANKWIRE][GETVIDEOBYID]",err)
			return SpankwireSingleVideo{}
		}
		JTCache.Put("spankwire-video-"+ID, b, spankwireCacheDuration)
	} else {
		json.Unmarshal(Cached.([]uint8), &result)
	}
	return result
}

func spankwireGetVideoEmbedCode(ID string) SpankwireEmbedCode {
	Cached := JTCache.Get("spankwire-embed-" + ID)
	if Cached == nil {
		timeout := time.Duration(spankwireAPITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(spankwireAPIURL+"?data=getVideoEmbedCode&output=json&video_id=%s", ID))
		if err != nil {
			log.Println("[SPANKWIRE][GETVIDEOEMBEDCODE]",err)
			return SpankwireEmbedCode{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result SpankwireEmbedCode
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[SPANKWIRE][GETVIDEOEMBEDCODE]",err)
		}
		return result
	}
	var result SpankwireEmbedCode
	json.Unmarshal(Cached.([]uint8), &result)
	return result		
}

// SpankwireSearch Calls spankwire search function and process result to get array of videos
func SpankwireSearch(search string) []JTVideo {
	videos := spankwireSearchVideos(search)
	result := []JTVideo{}
	for _, data := range videos["videos"].([]interface{}) {
		// Construct video object
		v := data.(map[string]interface{})["video"]
		video := JTVideo{}
		video.ID = fmt.Sprintf("%s", v.(map[string]interface{})["id"])
		if _, err := strconv.Atoi(video.ID); err != nil {
			video.ID = fmt.Sprintf("%.0f", v.(map[string]interface{})["id"])
		}
		video.Provider = "spankwire"
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
		/*tags := v.(map[string]interface{})["tags"] // Tags seem to be broken
		for _, tag := range tags.([]interface{}) {
			video.Tags = append(video.Tags, fmt.Sprintf("%s", tag.(map[string]interface{})["tag_name"]))
		}*/ 

		thumbs := v.(map[string]interface{})["thumbs"]
		for _, thumb := range thumbs.([]interface{}) {
			video.Thumbs = append(video.Thumbs, fmt.Sprintf("%s", thumb.(map[string]interface{})["src"]))
		}

		result = append(result, video)
	}

	return result
}

func spankwireSearchVideos(search string) KeezmoviesSearchResult {
	Cached := JTCache.Get("spankwire-search-" + search)
	if Cached == nil {
		timeout := time.Duration(spankwireAPITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(spankwireAPIURL+"?data=searchVideos&output=json&search=%s&thumbsize=small", url.QueryEscape(search)))

		if err != nil {
			log.Println("[SPANKWIRE][SEARCHVIDEOS]", err)
			return KeezmoviesSearchResult{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result KeezmoviesSearchResult
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[SPANKWIRE][SEARCHVIDEOS]", err)
			return KeezmoviesSearchResult{}
		}
		JTCache.Put("spankwire-search-"+search, b, spankwireCacheDuration)
		return result
	}
	var result KeezmoviesSearchResult
	json.Unmarshal(Cached.([]uint8), &result)
	return result
}