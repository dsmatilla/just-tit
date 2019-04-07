package main

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"io/ioutil"
	"net/http"
	"strings"
)

func imageProxy(image string) events.APIGatewayProxyResponse {
	aux := strings.Split(image, ".")
	str, _ := base64.StdEncoding.DecodeString(aux[0])

	response, _ := http.Get(fmt.Sprintf("%s", str))
	body, _ := ioutil.ReadAll(response.Body)
	c := base64.StdEncoding.EncodeToString(body)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "image/jpeg",
		},
		Body:            c,
		IsBase64Encoded: true,
	}
}
