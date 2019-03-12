package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dsmatilla/pornhub"
	"github.com/dsmatilla/redtube"
	"github.com/dsmatilla/tube8"
	"github.com/dsmatilla/youporn"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

const BaseDomain = "https://just-tit.com"

var AllowedDomains = []string{
	"just-tit.com",
	"dev.just-tit.com",
}

type searchResult struct {
	Pornhub pornhub.PornhubSearchResult
	Redtube redtube.RedtubeSearchResult
	Tube8 tube8.Tube8SearchResult
	Youporn youporn.YoupornSearchResult
}

var waitGroup sync.WaitGroup

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

func doSearch(search string) searchResult {
	waitGroup.Add(4)

	PornhubChannel := make(chan pornhub.PornhubSearchResult)
	RedtubeChannel := make(chan redtube.RedtubeSearchResult)
	Tube8Channel :=	make(chan tube8.Tube8SearchResult)
	YoupornChannel := make(chan youporn.YoupornSearchResult)

	go searchPornhub(search, PornhubChannel)
	go searchRedtube(search, RedtubeChannel)
	go searchTube8(search, Tube8Channel)
	go searchYouporn(search, YoupornChannel)

	result := searchResult{<-PornhubChannel, <-RedtubeChannel, <-Tube8Channel, <-YoupornChannel }

	waitGroup.Wait()

	return result
}

func search(search string) events.APIGatewayProxyResponse {
	headers := map[string]string{
		"Content-Type":"text/html; charset=utf-8",
		"Cache-Control":"max-age=31536000",
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
		PageTitle string
		Search string
		PageMetaDesc string
	}{
		PageTitle: fmt.Sprintf("Search results for %s", search),
		Search: search,
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
		Headers: headers,
		Body: body,
	}
}

func frontpage() (events.APIGatewayProxyResponse) {
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
		querystring = strings.Replace(querystring,"%20","_",-1)
		querystring = strings.Replace(querystring," ","_",-1)
		location := fmt.Sprintf("/%s.html", querystring)
		querystring = url.PathEscape(querystring)
		response.StatusCode = 301
		response.Headers = map[string]string{
			"Content-Type":"text/html",
			"Location":location,
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
			log.Fatal(err)
		}

		contentType := http.DetectContentType(fileContentBuffer)

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
