package main

import (
	_ "embed"
	"encoding/json"
	mapset "github.com/deckarep/golang-set/v2"
	"strings"
)

//go:embed resource/blacklist.json
var fileByte []byte

func filter(data []Data) ([]Data, error) {
	var blacklist Blacklist
	if err := json.Unmarshal(fileByte, &blacklist); err != nil {
		return nil, err
	}

	authors := mapset.NewThreadUnsafeSet(blacklist.Authors...)
	categories := mapset.NewThreadUnsafeSet(blacklist.Categories...)

	var ret []Data
Loop:
	for _, datum := range data {
		if authors.Contains(datum.Owner.Name) {
			continue
		}

		if categories.Contains(datum.Tname) {
			continue
		}

		for _, keyword := range blacklist.Keywords {
			if strings.Contains(datum.Title, keyword) || strings.Contains(datum.Desc, keyword) {
				continue Loop
			}
		}

		ret = append(ret, datum)
	}
	return ret, nil
}

type Blacklist struct {
	Authors    []string
	Categories []string
	Keywords   []string
}
