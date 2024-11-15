package crawler

import (
	"bytes"
	"encoding/base64"
	"html"
	"regexp"
	"strings"

	"github.com/nitezs/pcgamedb/constant"
	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/model"
	"github.com/nitezs/pcgamedb/utils"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type FreeGOGCrawler struct {
	logger *zap.Logger
}

// Deprecated: Unable to get through cloudflare
func NewFreeGOGCrawler(logger *zap.Logger) *FreeGOGCrawler {
	return &FreeGOGCrawler{
		logger: logger,
	}
}

func (c *FreeGOGCrawler) Crawl(num int) ([]*model.GameDownload, error) {
	count := 0
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.FreeGOGListURL,
	})
	if err != nil {
		c.logger.Error("Failed to fetch", zap.Error(err))
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		c.logger.Error("Failed to parse HTML", zap.Error(err))
		return nil, err
	}

	urls := []string{}
	updateFlags := []string{} //rawName+link
	doc.Find(".items-outer li a").Each(func(i int, s *goquery.Selection) {
		urls = append(urls, s.AttrOr("href", ""))
		updateFlags = append(updateFlags, s.Text()+s.AttrOr("href", ""))
	})

	res := []*model.GameDownload{}
	for i, u := range urls {
		if count == num {
			break
		}
		if db.IsFreeGOGCrawled(updateFlags[i]) {
			continue
		}
		c.logger.Info("Crawling", zap.String("URL", u))
		item, err := c.CrawlByUrl(u)
		if err != nil {
			c.logger.Warn("Failed to crawl", zap.Error(err), zap.String("URL", u))
			continue
		}
		item.UpdateFlag = updateFlags[i]
		err = db.SaveGameDownload(item)
		if err != nil {
			c.logger.Warn("Failed to save", zap.Error(err))
			continue
		}
		res = append(res, item)
		count++
		info, err := OrganizeGameDownload(item)
		if err != nil {
			c.logger.Warn("Failed to organize", zap.Error(err), zap.String("URL", u))
			continue
		}
		err = db.SaveGameInfo(info)
		if err != nil {
			c.logger.Warn("Failed to save", zap.Error(err), zap.String("URL", u))
			continue
		}
	}
	return res, nil
}

func (c *FreeGOGCrawler) CrawlByUrl(url string) (*model.GameDownload, error) {
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: url,
	})
	if err != nil {
		return nil, err
	}
	item, err := db.GetGameDownloadByUrl(url)
	if err != nil {
		return nil, err
	}
	item.Url = url
	rawTitleRegex := regexp.MustCompile(`(?i)<h1 class="entry-title">(.*?)</h1>`)
	rawTitleRegexRes := rawTitleRegex.FindStringSubmatch(string(resp.Data))
	rawName := ""
	if len(rawTitleRegexRes) > 1 {
		rawName = html.UnescapeString(rawTitleRegexRes[1])
		item.RawName = strings.Replace(rawName, "–", "-", -1)
	} else {
		return nil, err
	}
	item.Name = FreeGOGFormatter(item.RawName)
	sizeRegex := regexp.MustCompile(`(?i)>Size:\s?(.*?)<`)
	sizeRegexRes := sizeRegex.FindStringSubmatch(string(resp.Data))
	if len(sizeRegexRes) > 1 {
		item.Size = sizeRegexRes[1]
	}
	magnetRegex := regexp.MustCompile(`<a class="download-btn" href="https://gdl.freegogpcgames.xyz/download-gen\.php\?url=(.*?)"`)
	magnetRegexRes := magnetRegex.FindStringSubmatch(string(resp.Data))
	if len(magnetRegexRes) > 1 {
		magnet, err := base64.StdEncoding.DecodeString(magnetRegexRes[1])
		if err != nil {
			return nil, err
		}
		item.Download = string(magnet)
	} else {
		return nil, err
	}
	item.Author = "FreeGOG"
	return item, nil
}

func (c *FreeGOGCrawler) CrawlAll() ([]*model.GameDownload, error) {
	return c.Crawl(-1)
}

var freeGOGRegexps = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\(.*\)`),
}

func FreeGOGFormatter(name string) string {
	for _, re := range freeGOGRegexps {
		name = re.ReplaceAllString(name, "")
	}

	reg1 := regexp.MustCompile(`(?i)v\d+(\.\d+)*`)
	if index := reg1.FindIndex([]byte(name)); index != nil {
		name = name[:index[0]]
	}
	if index := strings.Index(name, "+"); index != -1 {
		name = name[:index]
	}

	reg2 := regexp.MustCompile(`(?i):\sgoty`)
	name = reg2.ReplaceAllString(name, ": Game Of The Year")

	return strings.TrimSpace(name)
}
