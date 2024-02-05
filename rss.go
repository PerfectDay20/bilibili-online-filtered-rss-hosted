package main

import (
	"bytes"
	"encoding/xml"
	"html/template"
	"strconv"
)

const title = "Filtered BiliBili online list"
const link = "https://www.bilibili.com/video/online.html"
const description = "A filtered BiliBili online list based on my blacklist"
const iconUrl = "https://www.bilibili.com/favicon.ico"

func encodeRss(bilibiliData *BilibiliData) string {
	rss := Rss{
		Version: "2.0",
		Channel: Channel{
			Title:       title,
			Link:        link,
			Description: description,
			Image:       Image{Title: title, Link: link, Url: iconUrl},
		},
	}

	rss.Channel.Item = createItem(bilibiliData)
	marshal, err := xml.MarshalIndent(&rss, "", "  ")
	if err != nil {
		return ""
	}

	return string(marshal)
}

func createItem(bilibiliData *BilibiliData) []Item {
	filteredData, _ := filter(bilibiliData.Data)

	items := make([]Item, len(filteredData))
	for i, data := range filteredData {
		items[i] = Item{
			Title:       data.Title,
			Link:        data.ShortLinkV2,
			Description: createItemDesc(&data),
			Guid:        data.ShortLinkV2,
		}
	}
	return items
}

// this is the content to show in the feed
// not just use the data.desc, add more info such as stats
func createItemDesc(data *Data) string {
	s := `<b>author:</b> {{.Owner.Name}}
	<p></p>
	<b>category:</b> {{.Tname}}
	<p></p>
	<b>desc:</b> {{.Desc}}
	<p></p>
	<b>view:</b> {{convertIntForHuman .Stat.View}}
	<p></p>
	<b>danmaku:</b> {{convertIntForHuman .Stat.Danmaku}}
	<p></p>
	<img style="width:100%" src="{{.Pic}}" width="500">`

	t := template.New("itemDesc")
	t = t.Funcs(template.FuncMap{"convertIntForHuman": convertIntForHuman})
	t = template.Must(t.Parse(s))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, data)
	if err != nil {
		return ""
	}
	return buf.String()
}

// Convert 1000 to 1k, 1000_000 to 1m, etc.
// Registered in template
func convertIntForHuman(i int) string {
	switch {
	case i < 1000:
		return strconv.Itoa(i)
	case i < 1000000:
		return strconv.Itoa(i/1000) + "k"
	default:
		return strconv.Itoa(i/1000_000) + "m"
	}
}

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Image       Image  `xml:"image"`
	Item        []Item `xml:"item"`
}

type Image struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Url   string `xml:"url"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Guid        string `xml:"guid"`
}
