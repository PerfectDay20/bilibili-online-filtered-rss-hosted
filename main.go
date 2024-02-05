package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io"
	"net/http"
)

func main() {
	lambda.Start(exec)
}

func exec() (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{StatusCode: 400}
	bytes, err := callApi()
	if err != nil {
		return response, err
	}
	bilibiliData, err := parseJson(bytes)
	if err != nil {
		return response, err
	}
	rssString := encodeRss(&bilibiliData)
	response.StatusCode = 200
	response.Headers = map[string]string{"content-type": "application/rss+xml"}
	response.Body = rssString
	return response, nil
}

func parseJson(bytes []byte) (BilibiliData, error) {
	var response BilibiliData

	if err := json.Unmarshal(bytes, &response); err != nil {
		return response, err
	}
	return response, nil
}

func callApi() ([]byte, error) {
	resp, err := http.Get("https://api.bilibili.com/x/web-interface/online/list")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil

}
