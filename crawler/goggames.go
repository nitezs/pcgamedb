package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/nitezs/pcgamedb/constant"
	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/model"
	"github.com/nitezs/pcgamedb/utils"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type GOGGamesCrawler struct {
	logger *zap.Logger
}

func NewGOGGamesCrawler(logger *zap.Logger) *GOGGamesCrawler {
	return &GOGGamesCrawler{
		logger: logger,
	}
}

func (c *GOGGamesCrawler) Name() string {
	return "GOGGamesCrawler"
}

func (c *GOGGamesCrawler) CrawlByUrl(url string) (*model.GameDownload, error) {
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
	name := strings.TrimSpace(doc.Find("#game-details>.container>h1").First().Text())
	magnetRegex := regexp.MustCompile(`magnet:\?[^"]*`)
	magnetRegexRes := magnetRegex.FindString(string(resp.Data))
	if magnetRegexRes == "" {
		return nil, errors.New("magnet not found")
	}
	sizeStrs := make([]string, 0)
	doc.Find(".container>.items-group").First().Find(".filesize").Each(func(i int, s *goquery.Selection) {
		sizeStrs = append(sizeStrs, s.Text())
	})
	size, err := utils.SubSizeStrings(sizeStrs)
	if err != nil {
		return nil, err
	}
	item, err := db.GetGameDownloadByUrl(url)
	if err != nil {
		return nil, err
	}
	item.Name = name
	item.RawName = name
	item.Download = magnetRegexRes
	item.Url = url
	item.Size = size
	item.Author = "GOGGames"
	return item, nil
}

func (c *GOGGamesCrawler) Crawl(page int) ([]*model.GameDownload, error) {
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: fmt.Sprintf(constant.GOGGamesURL, page),
	})
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		return nil, err
	}
	urls := make([]string, 0)
	doc.Find(".game-blocks>a").Each(func(i int, s *goquery.Selection) {
		u, exist := s.Attr("href")
		if !exist {
			return
		}
		urls = append(urls, fmt.Sprintf("%s%s", constant.GOGGamesBaseURL, u))
	})
	res := make([]*model.GameDownload, 0)
	for _, u := range urls {
		c.logger.Info("Crawling", zap.String("URL", u))
		item, err := c.CrawlByUrl(u)
		if err != nil {
			c.logger.Warn("Failed to crawl", zap.Error(err), zap.String("URL", u))
			continue
		}
		if err := db.SaveGameDownload(item); err != nil {
			c.logger.Warn("Failed to save", zap.Error(err), zap.String("URL", u))
			continue
		}
		res = append(res, item)
		info, err := OrganizeGameDownload(item)
		if err != nil {
			c.logger.Warn("Failed to organize", zap.Error(err), zap.String("URL", u))
			continue
		}
		if err := db.SaveGameInfo(info); err != nil {
			c.logger.Warn("Failed to save", zap.Error(err), zap.String("URL", u))
			continue
		}
	}
	return res, nil
}

func (c *GOGGamesCrawler) CrawlMulti(pages []int) ([]*model.GameDownload, error) {
	res := make([]*model.GameDownload, 0)
	for _, page := range pages {
		items, err := c.Crawl(page)
		if err != nil {
			return nil, err
		}
		res = append(res, items...)
	}
	return res, nil
}

func (c *GOGGamesCrawler) CrawlAll() ([]*model.GameDownload, error) {
	totalPageNum, err := c.GetTotalPageNum()
	if err != nil {
		return nil, err
	}
	var res []*model.GameDownload
	for i := 1; i <= totalPageNum; i++ {
		items, err := c.Crawl(i)
		if err != nil {
			return nil, err
		}
		res = append(res, items...)
	}
	return res, nil
}

func (c *GOGGamesCrawler) GetTotalPageNum() (int, error) {
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: fmt.Sprintf(constant.GOGGamesURL, 1),
	})
	if err != nil {
		return 0, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		return 0, err
	}
	btns := doc.Find(".pagination>.btn")
	return strconv.Atoi(strings.TrimSpace(btns.Eq(btns.Length() - 2).Text()))
}
