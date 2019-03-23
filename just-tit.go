package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/dsmatilla/pornhub"
	"github.com/dsmatilla/redtube"
	"github.com/dsmatilla/tube8"
	"github.com/dsmatilla/youporn"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const BaseDomain = "https://just-tit.com"

const dynamodbRegion = "eu-west-1"
const dynamodbTable = "JustTit"
const secondsToCache = 60 * 60 * 24

type JustTitCache struct {
	ID         string  `json:"id"`
	Result	   string  `json:"result"`
	Timestamp  int64 `json:"timestamp"`
}

var AllowedDomains = []string{
	"just-tit.com",
	"dev.just-tit.com",
}

type searchResult struct {
	Pornhub pornhub.PornhubSearchResult
	Redtube redtube.RedtubeSearchResult
	Tube8   tube8.Tube8SearchResult
	Youporn youporn.YoupornSearchResult
}

var waitGroup sync.WaitGroup

func getFromDB(ID string) JustTitCache {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(dynamodbRegion)},
	)
	svc := dynamodb.New(sess)

	result, _ := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(dynamodbTable),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(ID),
			},
		},
	})

	cache := JustTitCache{}
	if result != nil {
		dynamodbattribute.UnmarshalMap(result.Item, &cache)
	}

	return cache
}

func putToDB(ID string, Result string) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(dynamodbRegion)},
	)
	svc := dynamodb.New(sess)

	cache := JustTitCache{ID, Result, time.Now().Unix() + secondsToCache}
	item, _ := dynamodbattribute.MarshalMap(cache)
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(dynamodbTable),
	}
	_, _ = svc.PutItem(input)
}

func pornhubGetVideoByID(videoID string) pornhub.PornhubSingleVideo {
	cachedElement := getFromDB("pornhub-video-"+videoID)
	if (JustTitCache{}) != cachedElement {
		var result pornhub.PornhubSingleVideo
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		video := pornhub.GetVideoByID(videoID)
		json, _ := json.Marshal(video)
		putToDB("pornhub-video-"+videoID, string(json))
		return video
	}
}

func pornhubGetVideoEmbedCode(videoID string) pornhub.PornhubEmbedCode {
	cachedElement := getFromDB("pornhub-embed-"+videoID)
	if (JustTitCache{}) != cachedElement {
		var result pornhub.PornhubEmbedCode
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		embed := pornhub.GetVideoEmbedCode(videoID)
		json, _ := json.Marshal(embed)
		putToDB("pornhub-embed-"+videoID, string(json))
		return embed
	}
}

func redtubeGetVideoByID(videoID string) redtube.RedtubeSingleVideo {
	cachedElement := getFromDB("redtube-video-"+videoID)
	if (JustTitCache{}) != cachedElement {
		var result redtube.RedtubeSingleVideo
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		video := redtube.GetVideoByID(videoID)
		json, _ := json.Marshal(video)
		putToDB("redtube-video-"+videoID, string(json))
		return video
	}
}

func redtubeGetVideoEmbedCode(videoID string) redtube.RedtubeEmbedCode {
	cachedElement := getFromDB("redtube-embed-"+videoID)
	if (JustTitCache{}) != cachedElement {
		var result redtube.RedtubeEmbedCode
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		embed := redtube.GetVideoEmbedCode(videoID)
		json, _ := json.Marshal(embed)
		putToDB("redtube-embed-"+videoID, string(json))
		return embed
	}
}

func tube8GetVideoByID(videoID string) tube8.Tube8SingleVideo {
	cachedElement := getFromDB("tube8-video-"+videoID)
	if (JustTitCache{}) != cachedElement {
		var result tube8.Tube8SingleVideo
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		video := tube8.GetVideoByID(videoID)
		json, _ := json.Marshal(video)
		putToDB("tube8-video-"+videoID, string(json))
		return video
	}
}

func tube8GetVideoEmbedCode(videoID string) tube8.Tube8EmbedCode {
	cachedElement := getFromDB("tube8-embed-"+videoID)
	if (JustTitCache{}) != cachedElement {
		var result tube8.Tube8EmbedCode
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		embed := tube8.GetVideoEmbedCode(videoID)
		json, _ := json.Marshal(embed)
		putToDB("tube8-embed-"+videoID, string(json))
		return embed
	}
}

func youpornGetVideoByID(videoID string) youporn.YoupornSingleVideo {
	cachedElement := getFromDB("youporn-video-"+videoID)
	if (JustTitCache{}) != cachedElement {
		var result youporn.YoupornSingleVideo
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		video := youporn.GetVideoByID(videoID)
		json, _ := json.Marshal(video)
		putToDB("youporn-video-"+videoID, string(json))
		return video
	}
}

func youpornGetVideoEmbedCode(videoID string) youporn.YoupornEmbedCode {
	cachedElement := getFromDB("youporn-embed-"+videoID)
	if (JustTitCache{}) != cachedElement {
		var result youporn.YoupornEmbedCode
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		embed := youporn.GetVideoEmbedCode(videoID)
		json, _ := json.Marshal(embed)
		putToDB("youporn-embed-"+videoID, string(json))
		return embed
	}
}


func searchPornhub(search string, c chan pornhub.PornhubSearchResult) {
	defer waitGroup.Done()
	c <- pornhub.SearchVideos(search)
	close(c)
}

func searchRedtube(search string, c chan redtube.RedtubeSearchResult) {
	defer waitGroup.Done()
	c <- redtube.SearchVideos(search)
	close(c)
}

func searchTube8(search string, c chan tube8.Tube8SearchResult) {
	defer waitGroup.Done()
	c <- tube8.SearchVideos(search)
	close(c)
}

func searchYouporn(search string, c chan youporn.YoupornSearchResult) {
	defer waitGroup.Done()
	c <- youporn.SearchVideos(search)
	close(c)
}

func singlevideo(provider string, videoID string, tp string) events.APIGatewayProxyResponse {
	headers := map[string]string{
		"Content-Type":  "text/html; charset=utf-8",
		"Cache-Control": "max-age=31536000",
	}

	pre, _ := template.ParseFiles("html/single/singlevideo_pre.html")
	post, _ := template.ParseFiles("html/single/singlevideo_post.html")

	playerpre, _ := template.ParseFiles("html/single/player_pre.html")
	playerpost, _ := template.ParseFiles("html/single/player_post.html")

	// Build result divs
	var buff bytes.Buffer
	var embed string

	replace := struct {
		PageTitle    string
		Search       string
		PageMetaDesc string
		Url 		 string
		Thumb		 string
		Domain		 string
		Width 		 string
		Height       string
	}{
		PageTitle:    "",
		Search:       "",
		PageMetaDesc: "",
		Url:		  "",
		Thumb:		  "",
		Domain:		  BaseDomain,
		Width: 		  "",
		Height:       "",
	}

	switch provider {
	case "pornhub":
		video := pornhubGetVideoByID(videoID)
		embed = pornhubGetVideoEmbedCode(videoID).Embed.Code
		embed = fmt.Sprintf("%+v", html.UnescapeString(embed))
		replace.PageTitle = fmt.Sprintf("%s", video.Video.Title)
		replace.PageMetaDesc = fmt.Sprintf("%s", video.Video.Title)
		replace.Thumb = fmt.Sprintf("%s", video.Video.Thumb)
		replace.Url = fmt.Sprintf(BaseDomain+"/pornhub/%s.html", videoID)
		replace.Width = "580"
		replace.Height = "360"
	case "redtube":
		video := redtubeGetVideoByID(videoID)
		embed = redtubeGetVideoEmbedCode(videoID).Embed.Code
		str, _ := base64.StdEncoding.DecodeString(embed)
		embed = fmt.Sprintf("<object><embed src=\"%+v\" /></object>", html.UnescapeString(string(str)))
		replace.PageTitle = fmt.Sprintf("%s", video.Video.Title)
		replace.PageMetaDesc = fmt.Sprintf("%s", video.Video.Title)
		replace.Thumb = fmt.Sprintf("%s", video.Video.Thumb)
		replace.Url = fmt.Sprintf(BaseDomain+"/redtube/%s.html", videoID)
		replace.Width = "320"
		replace.Height = "180"
	case "tube8":
		video := tube8GetVideoByID(videoID)
		embed = tube8GetVideoEmbedCode(videoID).EmbedCode.Code
		embed = strings.Replace(embed, "![CDATA[", "", -1)
		embed = strings.Replace(embed, "]]", "", -1)
		str, _ := base64.StdEncoding.DecodeString(embed)
		embed = fmt.Sprintf("%+v", html.UnescapeString(string(str)))
		replace.PageTitle = fmt.Sprintf("%s", video.Videos.Title)
		replace.PageMetaDesc = fmt.Sprintf("%s", video.Videos.Title)
		replace.Thumb = fmt.Sprintf("%s", video.Videos.Thumbs.Thumb[0].Thumb)
		replace.Url = fmt.Sprintf(BaseDomain+"/tube8/%s.html", videoID)
		replace.Width = "628"
		replace.Height = "362"
	case "youporn":
		video := youpornGetVideoByID(videoID)
		embed = youpornGetVideoEmbedCode(videoID).Embed.Code
		embed = fmt.Sprintf("%+v", html.UnescapeString(embed))
		replace.PageTitle = fmt.Sprintf("%s", video.Video.Title)
		replace.PageMetaDesc = fmt.Sprintf("%s", video.Video.Title)
		replace.Thumb = fmt.Sprintf("%s", video.Video.Thumb)
		replace.Url = fmt.Sprintf(BaseDomain+"/youporn/%s.html", videoID)
		replace.Width = "628"
		replace.Height = "501"
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 301,
			Headers:    map[string]string{"Location": "/"},
			Body:       "",
		}
	}

	if tp == "true" {
		playerpre.Execute(&buff, replace)
		buff.WriteString(embed)
		playerpost.Execute(&buff, nil)
	} else {
		pre.Execute(&buff, replace)
		buff.WriteString(embed)
		post.Execute(&buff, nil)
	}
	body := buff.String()

	if len(embed) == 0 {
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

func doSearch(search string) searchResult {
		waitGroup.Add(4)

		PornhubChannel := make(chan pornhub.PornhubSearchResult)
		RedtubeChannel := make(chan redtube.RedtubeSearchResult)
		Tube8Channel := make(chan tube8.Tube8SearchResult)
		YoupornChannel := make(chan youporn.YoupornSearchResult)

		go searchPornhub(search, PornhubChannel)
		go searchRedtube(search, RedtubeChannel)
		go searchTube8(search, Tube8Channel)
		go searchYouporn(search, YoupornChannel)

		result := searchResult{<-PornhubChannel, <-RedtubeChannel, <-Tube8Channel, <-YoupornChannel}

		waitGroup.Wait()

		return result
}

func search(search string) events.APIGatewayProxyResponse {
	headers := map[string]string{
		"Content-Type":  "text/html; charset=utf-8",
		"Cache-Control": "max-age=31536000",
	}

	search = strings.Replace(search, ".html", "", -1)
	search = strings.Replace(search, "%20", " ", -1)
	search = strings.Replace(search, "_", " ", -1)
	result := doSearch(search)

	// Build HTML from template
	pre, _ := template.ParseFiles("html/search/searchresult_pre.html")
	divPornhub, _ := template.ParseFiles("html/search/tile_pornhub.html")
	divRedtube, _ := template.ParseFiles("html/search/tile_redtube.html")
	divTube8, _ := template.ParseFiles("html/search/tile_tube8.html")
	divYouporn, _ := template.ParseFiles("html/search/tile_youporn.html")
	post, _ := template.ParseFiles("html/search/searchresult_post.html")

	// Build result divs
	var buff bytes.Buffer
	replace := struct {
		PageTitle    string
		Search       string
		PageMetaDesc string
	}{
		PageTitle:    fmt.Sprintf("Search results for %s", search),
		Search:       search,
		PageMetaDesc: fmt.Sprintf("Search results for %s", search),
	}
	pre.Execute(&buff, replace)
	divPornhub.Execute(&buff, result.Pornhub.Videos)
	divRedtube.Execute(&buff, result.Redtube.Videos)
	divTube8.Execute(&buff, result.Tube8.Videos.Video)
	divYouporn.Execute(&buff, result.Youporn.Videos)
	post.Execute(&buff, nil)
	body := buff.String()

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       body,
	}
}

func frontpage() events.APIGatewayProxyResponse {
	headers := map[string]string{
		"Content-Type":  "text/html; charset=utf-8",
		"Cache-Control": "max-age=31536000",
	}

	// Build HTML from template
	html, _ := template.ParseFiles("html/frontpage/frontpage.html")

	replace := struct {
		PageTitle    string
		PageMetaDesc string
	}{
		PageTitle:    "Just Tit",
		PageMetaDesc: "The most optimized adult video search engine",
	}

	// Build result divs
	var buff bytes.Buffer
	html.Execute(&buff, replace)
	body := buff.String()

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       body,
	}
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	response := events.APIGatewayProxyResponse{}

	// Check if domains is allowed, redirect to base domain otherwise
	mustRedirect := true
	for _, domain := range AllowedDomains {
		if domain == request.Headers["Host"] {
			mustRedirect = false
		}
	}
	if mustRedirect {
		response.StatusCode = 301
		response.Headers = map[string]string{
			"Content-Type": "text/html",
			"Location":     BaseDomain + request.Path,
		}
		response.Body = ""
		return response, nil
	}

	if len(request.QueryStringParameters["s"]) > 0 {
		querystring := request.QueryStringParameters["s"]
		querystring = strings.Replace(querystring, "%20", "_", -1)
		querystring = strings.Replace(querystring, " ", "_", -1)
		location := fmt.Sprintf("/%s.html", querystring)
		querystring = url.PathEscape(querystring)
		response.StatusCode = 301
		response.Headers = map[string]string{
			"Content-Type": "text/html",
			"Location":     location,
		}
		return response, nil
	}

	// Frontpage
	if request.Path == "/" {
		response = frontpage()
		return response, nil
	}

	// If requested file exists in files/ , serve the file content with right mime type
	if _, err := os.Stat("./files" + request.Path); err == nil {
		fileContentBuffer, err := ioutil.ReadFile("./files" + request.Path)
		if err != nil {
			log.Println(err)
		}

		contentType := http.DetectContentType(fileContentBuffer)

		// Javascript
		aux := strings.Split(request.Path, ".")
		if(aux[len(aux) - 1]) == "js" {
			contentType = "application/javascript"
		}


		response.StatusCode = 200
		response.Headers = map[string]string{
			"Content-Type": contentType,
		}
		response.Body = string(fileContentBuffer)
		return response, nil
	}

	str := strings.Split(request.Path, "/")

	if len(str) == 2 {
		return search(str[1]), nil
	}

	if len(str) == 3 {
		provider := str[1]
		videoID := strings.Replace(str[2], ".html", "", -1)
		return singlevideo(provider, videoID, request.QueryStringParameters["tp"]), nil
	}

	response.StatusCode = 404
	response.Headers = map[string]string{
		"Content-Type": "text/html",
	}
	response.Body = "NOT FOUND"
	return response, nil
}

func main() {
	lambda.Start(handleRequest)
}
