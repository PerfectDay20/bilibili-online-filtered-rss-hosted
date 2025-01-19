package main

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func init() {
	// set slogger
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", exec)
	err := http.ListenAndServe(":8080", mux)
	slog.Error(err.Error())
}

var (
	cacheTime    = time.Time{}
	cacheContent = ""
)

func exec(w http.ResponseWriter, r *http.Request) {
	needUpdate := time.Since(cacheTime) > 10*time.Minute
	body := ""
	if needUpdate {
		slog.Info("Need update, call bilibili http api now")
		bytes, err := callApi()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		bilibiliData, err := parseJson(bytes)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		body = encodeRss(&bilibiliData)

		// update cache time and content
		cacheTime = time.Now()
		cacheContent = body
	} else {
		slog.Info("Return cached record")
		body = cacheContent
	}

	w.Header().Set("content-type", "application/rss+xml")
	w.Write([]byte(body))
}

func callApi() ([]byte, error) {
	startTime := time.Now()
	resp, err := http.Get("https://api.bilibili.com/x/web-interface/online/list")
	durationMs := (time.Since(startTime)).Milliseconds()
	slog.Info("call http api", "time_cost_ms", durationMs)

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
