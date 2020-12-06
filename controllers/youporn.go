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

const youpornAPIURL = "http://www.youporn.com/api/webmasters/"
const youpornAPITimeout = 3
const youpornCacheDuration = time.Minute * 5

// YoupornSearchResult type for youporn search result
type YoupornSearchResult map[string]interface{}

// YoupornSingleVideo type for youporn video result
type YoupornSingleVideo map[string]interface{}

// YoupornEmbedCode type for youporn embed code
type YoupornEmbedCode map[string]interface{}

// YoupornController Beego Controler
type YoupornController struct {
	beego.Controller
}

// Get Youporn Video controller
func (c *YoupornController) Get() {
    // Get videoID from URL
	aux := strings.Replace(c.Ctx.Request.URL.Path, ".html", "", -1)
	str := strings.Split(aux, "/")
	videoID := str[2]

    // Build redirect URL in case the API fails
	redirect := "https://www.youporn.com/watch/" + videoID + "/title/?utm_source=just-tit.com&utm_medium=embed&utm_campaign=hubtraffic_dsmatilla"

	// Get base domain from URL
	BaseDomain := "https://" + c.Controller.Ctx.Input.Domain()
	if c.Controller.Ctx.Input.Port() != 80 && c.Controller.Ctx.Input.Port() != 443 {
		BaseDomain += fmt.Sprintf("%s%d", ":", c.Controller.Ctx.Input.Port())
	}


	// Call the API and 307 redirect to fallback URL if something is not right
	data := youpornGetVideoByID(videoID)
	_, ok := data["video"]
	if !ok {
		c.Redirect(redirect, 307)
		return
	}

	// Get Embed Code from API
	embedcode := youpornGetVideoEmbedCode(videoID)
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
	video.Provider = "youporn"
	video.Domain = template.URL(BaseDomain)
	video.Title = fmt.Sprintf("%s", v["title"])
	video.Description = fmt.Sprintf("%s", v["title"])
	video.Thumb = fmt.Sprintf("%s", v["thumb"])
	video.Embed = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(embed["code"].(string))))
	video.URL = template.URL(fmt.Sprintf(BaseDomain+"/youporn/%s.html", videoID))
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
	for _, pornstars := range v["pornstars"].([]interface{}) {
		video.Pornstars = append(video.Pornstars, fmt.Sprintf("%s", pornstars.(map[string]interface{})["pornstar_name"]))
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

func youpornGetVideoByID(ID string) YoupornSingleVideo {
	Cached := JTCache.Get("youporn-video-"+ID)
	var result YoupornSingleVideo
	if Cached == nil {
		timeout := time.Duration(youpornAPITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(youpornAPIURL+"video_by_id/?video_id=%s", ID))
		if err != nil {
			log.Println("[YOUPORN][GETVIDEOBYID]",err)
			return YoupornSingleVideo{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[YOUPORN][GETVIDEOBYID]",err)
			return YoupornSingleVideo{}
		}
		JTCache.Put("youporn-video-"+ID, b, youpornCacheDuration)
	} else {
		json.Unmarshal(Cached.([]uint8), &result)
	}

	return result
}

func youpornGetVideoEmbedCode(ID string) YoupornEmbedCode {
	Cached := JTCache.Get("youporn-embed-"+ID)
	if Cached == nil {
		timeout := time.Duration(youpornAPITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(youpornAPIURL+"video_embed_code/?video_id=%s", ID))
		if err != nil {
			log.Println("[YOUPORN][GETVIDEOEMBEDCODE]",err)
			return YoupornEmbedCode{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result YoupornEmbedCode
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[YOUPORN][GETVIDEOEMBEDCODE]",err)
			return YoupornEmbedCode{}
		}
		JTCache.Put("youporn-embed-"+ID, b, youpornCacheDuration)
		return result
	}
	var result YoupornEmbedCode
	json.Unmarshal(Cached.([]uint8), &result)
	return result
}

// YoupornSearch Calls porhub search function and process result to get array of videos
func YoupornSearch(search string) []JTVideo {
	videos := youpornSearchVideos(search)
	result := []JTVideo{}
	if videos["video"] != nil {		
		for _, data := range videos["video"].([]interface{}) {
			// Construct video object
			//v := data.(map[string]interface{})["video"]
			v := data
			video := JTVideo{}
			video.ID = fmt.Sprintf("%s", v.(map[string]interface{})["video_id"])
			video.Provider = "youporn"
			video.Title = fmt.Sprintf("%s", v.(map[string]interface{})["title"])
			video.Description = fmt.Sprintf("%s", v.(map[string]interface{})["title"])
			video.Thumb = fmt.Sprintf("%s", v.(map[string]interface{})["thumb"])
			video.Width = fmt.Sprintf("%s", v.(map[string]interface{})["width"])
			video.Height = fmt.Sprintf("%s", v.(map[string]interface{})["height"])
			video.Duration = fmt.Sprintf("%s", v.(map[string]interface{})["duration"])
			video.Views = fmt.Sprintf("%s", v.(map[string]interface{})["views"])
			video.Rating = fmt.Sprintf("%s", v.(map[string]interface{})["rating"])
			video.Ratings = fmt.Sprintf("%s", v.(map[string]interface{})["ratings"])
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

func youpornSearchVideos(search string) YoupornSearchResult {
	Cached := JTCache.Get("youporn-search-"+search)
	if Cached == nil {
		timeout := time.Duration(youpornAPITimeout * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(fmt.Sprintf(youpornAPIURL+"search?search=%s&thumbsize=all", url.QueryEscape(search)))
		if err != nil {
			log.Println("[YOUPORN][SEARCHVIDEOS]",err)
			return YoupornSearchResult{}
		}
		b, _ := ioutil.ReadAll(resp.Body)
		var result YoupornSearchResult
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Println("[YOUPORN][SEARCHVIDEOS]",err)
			return YoupornSearchResult{}
		}
		JTCache.Put("youporn-search-"+search, b, youpornCacheDuration)
		return result
	}
	var result YoupornSearchResult
	json.Unmarshal(Cached.([]uint8), &result)
	return result
}