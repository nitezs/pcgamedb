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

type Formatter func(string) string

type s1337xCrawler struct {
	source    string
	formatter Formatter
	logger    *zap.Logger
}

func New1337xCrawler(source string, formatter Formatter, logger *zap.Logger) *s1337xCrawler {
	return &s1337xCrawler{
		source:    source,
		formatter: formatter,
		logger:    logger,
	}
}

func (c *s1337xCrawler) Crawl(page int) ([]*model.GameDownload, error) {
	var resp *utils.FetchResponse
	var doc *goquery.Document
	var err error
	requestUrl := fmt.Sprintf("%s/%s/%d/", constant.C1337xBaseURL, c.source, page)
	resp, err = utils.Fetch(utils.FetchConfig{
		Url: requestUrl,
	})
	if err != nil {
		return nil, err
	}
	doc, err = goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		return nil, err
	}
	trSelection := doc.Find("tbody>tr")
	urls := []string{}
	trSelection.Each(func(i int, trNode *goquery.Selection) {
		nameSelection := trNode.Find(".name").First()
		if aNode := nameSelection.Find("a").Eq(1); aNode.Length() > 0 {
			url, _ := aNode.Attr("href")
			urls = append(urls, url)
		}
	})
	var res []*model.GameDownload
	for _, u := range urls {
		u = fmt.Sprintf("%s%s", constant.C1337xBaseURL, u)
		if db.IsGameCrawledByURL(u) {
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
			c.logger.Warn("Failed to save", zap.Error(err), zap.String("URL", u))
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

func (c *s1337xCrawler) CrawlByUrl(url string) (*model.GameDownload, error) {
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: url,
	})
	if err != nil {
		return nil, err
	}
	var item = &model.GameDownload{}
	item.Url = url
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		return nil, err
	}
	selection := doc.Find(".torrent-detail-page ul.list>li")
	info := make(map[string]string)
	selection.Each(func(i int, item *goquery.Selection) {
		info[strings.TrimSpace(item.Find("strong").Text())] = strings.TrimSpace(item.Find("span").Text())
	})
	magnetRegex := regexp.MustCompile(`magnet:\?[^"]*`)
	magnetRegexRes := magnetRegex.FindStringSubmatch(string(resp.Data))
	item.Size = info["Total size"]
	item.RawName = doc.Find("title").Text()
	item.RawName = strings.Replace(item.RawName, "Download ", "", 1)
	item.RawName = strings.TrimSpace(strings.Replace(item.RawName, "Torrent | 1337x", " ", 1))
	item.Name = c.formatter(item.RawName)
	item.Download = magnetRegexRes[0]
	item.Author = strings.Replace(c.source, "-torrents", "", -1)
	return item, nil
}

func (c *s1337xCrawler) CrawlMulti(pages []int) (res []*model.GameDownload, err error) {
	var items []*model.GameDownload
	totalPageNum, err := c.GetTotalPageNum()
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		if page > totalPageNum {
			continue
		}
		items, err = c.Crawl(page)
		res = append(res, items...)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (c *s1337xCrawler) CrawlAll() (res []*model.GameDownload, err error) {
	totalPageNum, err := c.GetTotalPageNum()
	if err != nil {
		return nil, err
	}
	var items []*model.GameDownload
	for i := 1; i <= totalPageNum; i++ {
		items, err = c.Crawl(i)
		res = append(res, items...)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (c *s1337xCrawler) GetTotalPageNum() (int, error) {
	var resp *utils.FetchResponse
	var doc *goquery.Document
	var err error

	requestUrl := fmt.Sprintf("%s/%s/%d/", constant.C1337xBaseURL, c.source, 1)
	resp, err = utils.Fetch(utils.FetchConfig{
		Url: requestUrl,
	})
	if err != nil {
		return 0, err
	}
	doc, _ = goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	selection := doc.Find(".last")
	pageStr, exist := selection.Find("a").Attr("href")
	if !exist {
		return 0, errors.New("total page num not found")
	}
	pageStr = strings.ReplaceAll(pageStr, c.source, "")
	pageStr = strings.ReplaceAll(pageStr, "/", "")
	totalPageNum, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, err
	}
	return totalPageNum, nil
}
