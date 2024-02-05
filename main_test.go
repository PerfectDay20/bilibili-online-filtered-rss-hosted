package main

import (
	_ "embed"
	"fmt"
	"os"
	"testing"
)

func TestLocal(t *testing.T) {
	bytes, err := os.ReadFile("resource/bilibili.json")
	if err != nil {
		return
	}
	bilibiliData, err := parseJson(bytes)
	if err != nil {
		return
	}
	rss := encodeRss(&bilibiliData)
	fmt.Println(rss)

}
