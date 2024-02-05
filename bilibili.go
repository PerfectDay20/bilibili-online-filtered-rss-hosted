package main

type BilibiliData struct {
	Code    int
	Message string
	Ttl     int
	Data    []Data
}

type Data struct {
	Tname       string
	Pic         string
	Title       string
	Desc        string
	Owner       Owner
	Stat        Stat
	ShortLinkV2 string `json:"short_link_v2"`
}

type Owner struct {
	Name string
}

type Stat struct {
	View     int
	Danmaku  int
	Reply    int
	Favorite int
	Coin     int
	Share    int
}
