package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/nitezs/pcgamedb/constant"
	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/model"
	"github.com/nitezs/pcgamedb/utils"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type SteamRIPCrawler struct {
	logger *zap.Logger
}

func NewSteamRIPCrawler(logger *zap.Logger) *SteamRIPCrawler {
	return &SteamRIPCrawler{
		logger: logger,
	}
}

func (c *SteamRIPCrawler) Name() string {
	return "SteamRIPCrawler"
}

func (c *SteamRIPCrawler) CrawlByUrl(url string) (*model.GameItem, error) {
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: url,
	})
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		return nil, err
	}
	item, err := db.GetGameItemByUrl(url)
	if err != nil {
		return nil, err
	}
	item.RawName = strings.TrimSpace(doc.Find(".entry-title").First().Text())
	item.Name = SteamRIPFormatter(item.RawName)
	item.Url = url
	item.Author = "SteamRIP"
	sizeRegex := regexp.MustCompile(`(?i)<li><strong>Game Size:\s?</strong>(.*?)</li>`)
	sizeRegexRes := sizeRegex.FindStringSubmatch(string(resp.Data))
	if len(sizeRegexRes) != 0 {
		item.Size = strings.TrimSpace(sizeRegexRes[1])
	} else {
		item.Size = "unknown"
	}
	megadbRegex := regexp.MustCompile(`(?i)(?:https?:)?(//megadb\.net/[^"]+)`)
	megadbRegexRes := megadbRegex.FindStringSubmatch(string(resp.Data))
	if len(megadbRegexRes) != 0 {
		item.Download = fmt.Sprintf("https:%s", megadbRegexRes[1])
	}
	if item.Download == "" {
		gofileRegex := regexp.MustCompile(`(?i)(?:https?:)?(//gofile\.io/d/[^"]+)`)
		gofileRegexRes := gofileRegex.FindStringSubmatch(string(resp.Data))
		if len(gofileRegexRes) != 0 {
			item.Download = fmt.Sprintf("https:%s", gofileRegexRes[1])
		}
	}
	if item.Download == "" {
		filecryptRegex := regexp.MustCompile(`(?i)(?:https?:)?(//filecrypt\.co/Container/[^"]+)`)
		filecryptRegexRes := filecryptRegex.FindStringSubmatch(string(resp.Data))
		if len(filecryptRegexRes) != 0 {
			item.Download = fmt.Sprintf("https:%s", filecryptRegexRes[1])
		}
	}
	if item.Download == "" {
		return nil, errors.New("Failed to find download link")
	}

	return item, nil
}

func (c *SteamRIPCrawler) Crawl(num int) ([]*model.GameItem, error) {
	count := 0
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.SteamRIPGameListURL,
	})
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		return nil, err
	}
	var items []*model.GameItem
	urls := []string{}
	updateFlags := []string{} // title
	doc.Find(".az-list-item>a").Each(func(i int, s *goquery.Selection) {
		u, exist := s.Attr("href")
		if !exist {
			return
		}
		urls = append(urls, fmt.Sprintf("%s%s", constant.SteamRIPBaseURL, u))
		updateFlags = append(updateFlags, s.Text())
	})
	for i, u := range urls {
		if count == num {
			break
		}
		if db.IsSteamRIPCrawled(updateFlags[i]) {
			continue
		}
		c.logger.Info("Crawling", zap.String("URL", u))
		item, err := c.CrawlByUrl(u)
		if err != nil {
			c.logger.Error("Failed to crawl", zap.Error(err), zap.String("URL", u))
			continue
		}
		item.UpdateFlag = updateFlags[i]
		if err := db.SaveGameItem(item); err != nil {
			c.logger.Error("Failed to save item", zap.Error(err))
			continue
		}
		items = append(items, item)
		count++
		info, err := OrganizeGameItem(item)
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
	return items, nil
}

func (c *SteamRIPCrawler) CrawlAll() ([]*model.GameItem, error) {
	return c.Crawl(-1)
}

func SteamRIPFormatter(name string) string {
	name = regexp.MustCompile(`\([^\)]+\)`).ReplaceAllString(name, "")
	name = strings.Replace(name, "Free Download", "", -1)
	name = strings.TrimSpace(name)
	return name
}
