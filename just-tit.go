package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var BaseDomain = os.Getenv("BaseDomain")
var Theme = os.Getenv("Theme")

type TemplateData struct {
	PageTitle    string
	Search       string
	PageMetaDesc string
	Url          string
	Thumb        string
	Domain       string
	Width        string
	Height       string
	Embed        template.HTML
	Result       searchResult
}

var TemplateFunctions = template.FuncMap{
	"ToImageProxy": func(url string) string {
		if os.Getenv("ImageProxy") == "yes" {
			aux := strings.Split(url, ".")
			ext := aux[len(aux)-1]
			if ext == "jpg" {
				return BaseDomain + "/images/" + base64.StdEncoding.EncodeToString([]byte(url)) + "." + ext
			} else {
				return url
			}
		} else {
			return url
		}
	},
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	response := events.APIGatewayProxyResponse{}

	u, err := url.Parse(BaseDomain)
	if err != nil {
		log.Println("[JUST-TIT][URL_PARSE]", err)
	}

	if u.Host != request.Headers["Host"] {
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
		if (aux[len(aux)-1]) == "js" {
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
		if str[1] == "images" {
			return imageProxy(str[2]), nil
		} else {
			provider := str[1]
			videoID := strings.Replace(str[2], ".html", "", -1)
			return singlevideo(provider, videoID, request.QueryStringParameters["tp"]), nil
		}
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
