package utils

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func SolveKeepLinks(url string) (string, error) {
	id := url[strings.LastIndex(url, "/")+1:]
	resp, err := Fetch(FetchConfig{
		Url: url,
		Cookies: map[string]string{
			fmt.Sprintf("flag[%s]", id): "1",
		},
	})
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		return "", err
	}
	return doc.Find(".livelbl a").Text(), nil
}
