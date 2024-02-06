package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	lambda.Start(exec)
}

func exec() (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{StatusCode: 400}

	// check dynamodb cache
	table := initTable()
	record, err := table.GetRecord()
	if err != nil {
		return response, err
	}
	needUpdate := false
	if len(record.Record) == 0 {
		log.Println("DynamoDB cache empty, need to update")
		needUpdate = true
	} else if time.Now().Sub(time.Unix(record.UpdateTimestamp, 0)) > 10*time.Minute {
		log.Println("Cached content expired, need to update")
		needUpdate = true
	}

	if needUpdate {
		log.Println("Call bilibili http api now")
		bytes, err := callApi()
		if err != nil {
			return response, err
		}
		bilibiliData, err := parseJson(bytes)
		if err != nil {
			return response, err
		}
		rssString := encodeRss(&bilibiliData)
		err = table.SetRecord(rssString)
		if err != nil {
			log.Printf("Set record failed, reason: %v\n", err)
		}
		response.Body = rssString
	} else {
		log.Println("Return cached record")
		response.Body = record.Record
	}

	response.StatusCode = 200
	response.Headers = map[string]string{"content-type": "application/rss+xml"}
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
