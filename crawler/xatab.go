package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"pcgamedb/constant"
	"pcgamedb/db"
	"pcgamedb/model"
	"pcgamedb/utils"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type XatabCrawler struct {
	logger *zap.Logger
}

func NewXatabCrawler(logger *zap.Logger) *XatabCrawler {
	return &XatabCrawler{
		logger: logger,
	}
}

func (c *XatabCrawler) Crawl(page int) ([]*model.GameDownload, error) {
	requestURL := fmt.Sprintf("%s/page/%v", constant.XatabBaseURL, page)
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: requestURL,
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
	updateFlags := []string{} // title
	doc.Find(".entry").Each(func(i int, s *goquery.Selection) {
		u, exist := s.Find(".entry__title.h2 a").Attr("href")
		if !exist {
			return
		}
		urls = append(urls, u)
		updateFlags = append(updateFlags, s.Find(".entry__title.h2 a").Text())
	})
	var res []*model.GameDownload
	for i, u := range urls {
		if db.IsXatabCrawled(updateFlags[i]) {
			continue
		}
		c.logger.Info("Crawling", zap.String("URL", u))
		item, err := c.CrawlByUrl(u)
		if err != nil {
			c.logger.Warn("Failed to crawl", zap.Error(err), zap.String("URL", u))
			continue
		}
		err = db.SaveGameDownload(item)
		if err != nil {
			c.logger.Warn("Failed to save", zap.Error(err))
			continue
		}
		res = append(res, item)
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

func (c *XatabCrawler) CrawlByUrl(url string) (*model.GameDownload, error) {
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
	item, err := db.GetGameDownloadByUrl(url)
	if err != nil {
		return nil, err
	}
	item.Url = url
	item.RawName = doc.Find(".inner-entry__title").First().Text()
	item.Name = XatabFormatter(item.RawName)
	item.Author = "Xatab"
	item.UpdateFlag = item.RawName
	downloadURL := doc.Find("#download>a").First().AttrOr("href", "")
	if downloadURL == "" {
		return nil, errors.New("Failed to find download URL")
	}
	resp, err = utils.Fetch(utils.FetchConfig{
		Headers: map[string]string{"Referer": url},
		Url:     downloadURL,
	})
	if err != nil {
		return nil, err
	}
	magnet, size, err := utils.ConvertTorrentToMagnet(resp.Data)
	if err != nil {
		return nil, err
	}
	item.Size = size
	item.Download = magnet
	return item, nil
}

func (c *XatabCrawler) CrawlMulti(pages []int) ([]*model.GameDownload, error) {
	totalPageNum, err := c.GetTotalPageNum()
	if err != nil {
		return nil, err
	}
	var res []*model.GameDownload
	for _, page := range pages {
		if page > totalPageNum {
			continue
		}
		items, err := c.Crawl(page)
		if err != nil {
			return nil, err
		}
		res = append(res, items...)
	}
	return res, nil
}

func (c *XatabCrawler) CrawlAll() ([]*model.GameDownload, error) {
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

func (c *XatabCrawler) GetTotalPageNum() (int, error) {
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.XatabBaseURL,
	})
	if err != nil {
		return 0, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		return 0, err
	}
	pageStr := doc.Find(".pagination>a").Last().Text()
	totalPageNum, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, err
	}
	return totalPageNum, nil
}

var xatabRegexps = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\sPC$`),
}

func XatabFormatter(name string) string {
	reg1 := regexp.MustCompile(`(?i)v(er)?\s?(\.)?\d+(\.\d+)*`)
	if index := reg1.FindIndex([]byte(name)); index != nil {
		name = name[:index[0]]
	}
	if index := strings.Index(name, "["); index != -1 {
		name = name[:index]
	}
	if index := strings.Index(name, "("); index != -1 {
		name = name[:index]
	}
	if index := strings.Index(name, "{"); index != -1 {
		name = name[:index]
	}
	if index := strings.Index(name, "+"); index != -1 {
		name = name[:index]
	}
	name = strings.TrimSpace(name)
	for _, re := range xatabRegexps {
		name = re.ReplaceAllString(name, "")
	}

	if index := strings.Index(name, "/"); index != -1 {
		names := strings.Split(name, "/")
		longestLength := 0
		longestName := ""
		for _, n := range names {
			if !utils.ContainsRussian(n) && len(n) > longestLength {
				longestLength = len(n)
				longestName = n
			}
		}
		name = longestName
	}

	return strings.TrimSpace(name)
}
