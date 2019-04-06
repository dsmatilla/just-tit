package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"html/template"
	"strings"
)

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
	web := template.Must(template.New("search").Funcs(TemplateFunctions).ParseFiles(
		"html/"+Theme+"/template.html",
		"html/"+Theme+"/search/container.html",
		"html/"+Theme+"/search/pornhub.html",
		"html/"+Theme+"/search/redtube.html",
		"html/"+Theme+"/search/tube8.html",
		"html/"+Theme+"/search/youporn.html",
	))

	// Build result divs
	var buff bytes.Buffer
	replace := TemplateData{
		PageTitle:    fmt.Sprintf("Search results for %s", search),
		Search:       search,
		PageMetaDesc: fmt.Sprintf("Search results for %s", search),
		Result:       result,
	}

	web.ExecuteTemplate(&buff, "layout", replace)

	body := buff.String()

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       body,
	}
}
