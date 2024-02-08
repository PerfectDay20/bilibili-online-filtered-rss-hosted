package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var table *TableBasics

func init() {
	// set slogger
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	// get cached rss content in dynamo
	t, err := initDynamoTable()
	if err != nil {
		slog.Error("Can't access dynamo table, abort now")
		os.Exit(2)
	}
	table = t
}

func main() {
	lambda.Start(exec)
}

func exec() (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{StatusCode: 400}

	// check dynamodb cache
	needUpdate := false
	record, err := table.GetRecord()
	if err != nil {
		return response, err
	}

	if len(record.Record) == 0 {
		slog.Info("DynamoDB cache empty, need to update")
		needUpdate = true
	} else if time.Now().Sub(time.Unix(record.UpdateTimestamp, 0)) > 10*time.Minute {
		slog.Info("Cached content expired, need to update")
		needUpdate = true
	}

	if needUpdate {
		slog.Info("Call bilibili http api now")
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
			slog.Error("Set record failed", "reason", err)
		}
		response.Body = rssString
	} else {
		slog.Info("Return cached record")
		response.Body = record.Record
	}

	response.StatusCode = 200
	response.Headers = map[string]string{"content-type": "application/rss+xml"}
	return response, nil
}

func callApi() ([]byte, error) {
	startTime := time.Now()
	resp, err := http.Get("https://api.bilibili.com/x/web-interface/online/list")
	durationMs := (time.Now().Sub(startTime)).Milliseconds()
	slog.Info("call http api time", "time", durationMs)

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
