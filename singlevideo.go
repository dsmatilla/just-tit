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
	headers := map[string]string{
		"Content-Type":  "text/html; charset=utf-8",
		"Cache-Control": "max-age=31536000",
	}

	var templateFile string
	if tp == "true" {
		templateFile = "html/"+Theme+"/player.html"
	} else {
		templateFile = "html/"+Theme+"/template.html"
	}

	web, _ := template.ParseFiles(
		templateFile,
		"html/"+Theme+"/video/container.html",
	)

	replace := TemplateData{
		Domain:       BaseDomain,
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
	case "tube8":
		video := tube8GetVideoByID(videoID)
		embed := tube8GetVideoEmbedCode(videoID).EmbedCode.Code
		embed = strings.Replace(embed, "![CDATA[", "", -1)
		embed = strings.Replace(embed, "]]", "", -1)
		str, _ := base64.StdEncoding.DecodeString(embed)
		replace.Embed = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(string(str))))
		replace.PageTitle = fmt.Sprintf("%s", video.Videos.Title)
		replace.PageMetaDesc = fmt.Sprintf("%s", video.Videos.Title)
		replace.Thumb = fmt.Sprintf("%s", video.Videos.Thumbs.Thumb[0].Thumb)
		replace.Url = fmt.Sprintf(BaseDomain+"/tube8/%s.html", videoID)
		replace.Width = "628"
		replace.Height = "362"
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
	case "xtube":
		video := xtubeGetVideoByID(videoID)
		replace.Embed = template.HTML(fmt.Sprintf("<object><embed src=\"%+v\" /></object>", video.EmbedCode))
		replace.PageTitle = fmt.Sprintf("%s", video.Title)
		replace.PageMetaDesc = fmt.Sprintf("%s", video.Description)
		replace.Thumb = fmt.Sprintf("%s", video.Thumb)
		replace.Url = fmt.Sprintf(BaseDomain+"/xtube/%s.html", videoID)
		replace.Width = "628"
		replace.Height = "501"
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
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 301,
			Headers:    map[string]string{"Location": "/"},
			Body:       "",
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
