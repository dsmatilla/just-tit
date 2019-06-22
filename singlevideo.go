package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"html"
	"html/template"
	"strings"
)

func singlevideo(provider string, videoID string, tp string) events.APIGatewayProxyResponse {
	var headers = map[string]string{}

	headers = map[string]string{
		"Content-Type":  "text/html; charset=utf-8",
		"Cache-Control": "max-age=31536000",
	}

	var templateFile string
	if tp == "true" {
		templateFile = "html/" + Theme + "/player.html"
	} else {
		templateFile = "html/" + Theme + "/template.html"
	}

	// Build HTML from template
	web := template.Must(template.New("singlevideo").Funcs(TemplateFunctions).ParseFiles(
		templateFile,
		"html/"+Theme+"/video/container.html",
		"html/"+Theme+"/search/container.html",
		"html/"+Theme+"/search/pornhub.html",
		"html/"+Theme+"/search/redtube.html",
		"html/"+Theme+"/search/tube8.html",
		"html/"+Theme+"/search/youporn.html",
	))

	replace := TemplateData{
		ID:		videoID,
		Domain: BaseDomain,
	}

	switch provider {
	case "pornhub":
		video := pornhubGetVideoByID(videoID)
		embed := pornhubGetVideoEmbedCode(videoID).Embed.Code
		replace.Embed = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(embed)))
		replace.PageTitle = fmt.Sprintf("%s", video.Video.Title)
		replace.PageMetaDesc = fmt.Sprintf("%s", video.Video.Title)
		replace.Thumb = fmt.Sprintf("%s", video.Video.Thumb)
		replace.Url = fmt.Sprintf(BaseDomain+"/pornhub/%s.html", videoID)
		replace.Width = "580"
		replace.Height = "360"
		replace.PornhubVideo = video
	case "redtube":
		video := redtubeGetVideoByID(videoID)
		embed := redtubeGetVideoEmbedCode(videoID).Embed.Code
		str, _ := base64.StdEncoding.DecodeString(embed)
		replace.Embed = template.HTML(fmt.Sprintf("<object><embed src=\"%+v\" /></object>", html.UnescapeString(string(str))))
		replace.PageTitle = fmt.Sprintf("%s", video.Video.Title)
		replace.PageMetaDesc = fmt.Sprintf("%s", video.Video.Title)
		replace.Thumb = fmt.Sprintf("%s", video.Video.Thumb)
		replace.Url = fmt.Sprintf(BaseDomain+"/redtube/%s.html", videoID)
		replace.Width = "320"
		replace.Height = "180"
		replace.RedtubeVideo = video
	case "tube8":
		video := tube8GetVideoByID(videoID)
		embed := tube8GetVideoEmbedCode(videoID).EmbedCode.Code
		embed = strings.Replace(embed, "![CDATA[", "", -1)
		embed = strings.Replace(embed, "]]", "", -1)
		str, _ := base64.StdEncoding.DecodeString(embed)
		replace.Embed = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(string(str))))
		replace.PageTitle = fmt.Sprintf("%s", video.Videos.Title)
		replace.PageMetaDesc = fmt.Sprintf("%s", video.Videos.Title)
		if len(video.Videos.Thumbs.Thumb) > 0 {
			replace.Thumb = fmt.Sprintf("%s", video.Videos.Thumbs.Thumb[0].Thumb)
		}
		replace.Url = fmt.Sprintf(BaseDomain+"/tube8/%s.html", videoID)
		replace.Width = "628"
		replace.Height = "362"
		replace.Tube8Video = video
	case "youporn":
		video := youpornGetVideoByID(videoID)
		embed := youpornGetVideoEmbedCode(videoID).Embed.Code
		replace.Embed = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(embed)))
		replace.PageTitle = fmt.Sprintf("%s", video.Video.Title)
		replace.PageMetaDesc = fmt.Sprintf("%s", video.Video.Title)
		replace.Thumb = fmt.Sprintf("%s", video.Video.Thumb)
		replace.Url = fmt.Sprintf(BaseDomain+"/youporn/%s.html", videoID)
		replace.Width = "628"
		replace.Height = "501"
		replace.YoupornVideo = video
	case "xtube":
		video := xtubeGetVideoByID(videoID)
		replace.Embed = template.HTML(fmt.Sprintf("<object><embed src=\"%+v\" /></object>", video.EmbedCode))
		replace.PageTitle = fmt.Sprintf("%s", video.Title)
		replace.PageMetaDesc = fmt.Sprintf("%s", video.Description)
		replace.Thumb = fmt.Sprintf("%s", video.Thumb)
		replace.Url = fmt.Sprintf(BaseDomain+"/xtube/%s.html", videoID)
		replace.Width = "628"
		replace.Height = "501"
		replace.XtubeVideo = video
	case "spankwire":
		video := spankwireGetVideoByID(videoID)
		embed := spankwireGetVideoEmbedCode(videoID).Embed.Code
		str, _ := base64.StdEncoding.DecodeString(embed)
		replace.Embed = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(string(str))))
		replace.PageTitle = fmt.Sprintf("%s", video.Video.Title)
		replace.PageMetaDesc = fmt.Sprintf("%s", video.Video.Title)
		replace.Thumb = fmt.Sprintf("%s", video.Video.Thumb)
		replace.Url = fmt.Sprintf(BaseDomain+"/spankwire/%s.html", videoID)
		replace.Width = "650"
		replace.Height = "550"
		replace.SpankwireVideo = video
	case "keezmovies":
		video := keezmoviesGetVideoByID(videoID)
		embed := keezmoviesGetVideoEmbedCode(videoID).Embed.Code
		str, _ := base64.StdEncoding.DecodeString(embed)
		replace.Embed = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(string(str))))
		replace.PageTitle = fmt.Sprintf("%s", video.Video.Title)
		replace.PageMetaDesc = fmt.Sprintf("%s", video.Video.Title)
		replace.Thumb = fmt.Sprintf("%s", video.Video.Thumb)
		replace.Url = fmt.Sprintf(BaseDomain+"/keezmovies/%s.html", videoID)
		replace.Width = "650"
		replace.Height = "550"
		replace.KeezmoviesVideo = video
	case "extremetube":
		video := extremetubeGetVideoByID(videoID)
		embed := extremetubeGetVideoEmbedCode(videoID).Embed.Code
		str, _ := base64.StdEncoding.DecodeString(embed)
		replace.Embed = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(string(str))))
		replace.PageTitle = fmt.Sprintf("%s", video.Video.Title)
		replace.PageMetaDesc = fmt.Sprintf("%s", video.Video.Title)
		replace.Thumb = fmt.Sprintf("%s", video.Video.Thumb)
		replace.Url = fmt.Sprintf(BaseDomain+"/extremetube/%s.html", videoID)
		replace.Width = "650"
		replace.Height = "550"
		replace.ExtremetubeVideo = video
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 301,
			Headers:    map[string]string{"Location": "/"},
			Body:       "",
		}
	}

	replace.Result = doSearch(replace.PageTitle)
	replace.Result.Flag = false


	for index, element := range replace.Result.Pornhub.Videos {
		if element.ID == videoID {
			replace.Result.Pornhub.Videos = append(replace.Result.Pornhub.Videos[:index], replace.Result.Pornhub.Videos[index+1:]...)
		}
	}

	for index, element := range replace.Result.Redtube.Videos {
		if element.Video.ID == videoID {
			replace.Result.Redtube.Videos = append(replace.Result.Redtube.Videos[:index], replace.Result.Redtube.Videos[index+1:]...)
		}
	}

	for index, element := range replace.Result.Tube8.Videos.Video {
		if element.ID == videoID {
			replace.Result.Tube8.Videos.Video = append(replace.Result.Tube8.Videos.Video[:index], replace.Result.Tube8.Videos.Video[index+1:]...)
		}
	}

	for index, element := range replace.Result.Youporn.Videos {
		if element.ID == videoID {
			replace.Result.Youporn.Videos = append(replace.Result.Youporn.Videos[:index], replace.Result.Youporn.Videos[index+1:]...)
		}
	}

	var buff bytes.Buffer
	web.ExecuteTemplate(&buff, "layout", replace)

	body := buff.String()

	if len(replace.PageTitle) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: 301,
			Headers: map[string]string{
				"Content-Type": "text/html",
				"Location":     BaseDomain,
			},
			Body: "",
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       body,
	}
}
