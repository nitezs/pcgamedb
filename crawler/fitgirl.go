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

type FitGirlCrawler struct {
	logger *zap.Logger
}

func NewFitGirlCrawler(logger *zap.Logger) *FitGirlCrawler {
	return &FitGirlCrawler{
		logger: logger,
	}
}

func (c *FitGirlCrawler) CrawlByUrl(url string) (*model.GameDownload, error) {
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
	titleElem := doc.Find("h3").First().Find("strong")
	if titleElem.Length() == 0 {
		return nil, errors.New("Failed to find title")
	}
	rawTitle := titleElem.Text()
	titleElem.Children().Remove()
	title := strings.TrimSpace(titleElem.Text())
	sizeRegex := regexp.MustCompile(`Repack Size: <strong>(.*?)</strong>`)
	sizeRegexRes := sizeRegex.FindStringSubmatch(string(resp.Data))
	if len(sizeRegexRes) == 0 {
		return nil, errors.New("Failed to find size")
	}
	size := sizeRegexRes[1]
	magnetRegex := regexp.MustCompile(`magnet:\?[^"]*`)
	magnetRegexRes := magnetRegex.FindStringSubmatch(string(resp.Data))
	if len(magnetRegexRes) == 0 {
		return nil, errors.New("Failed to find magnet")
	}
	magnet := magnetRegexRes[0]
	item, err := db.GetGameDownloadByUrl(url)
	if err != nil {
		return nil, err
	}
	item.Name = strings.TrimSpace(title)
	item.RawName = rawTitle
	item.Url = url
	item.Size = size
	item.Author = "FitGirl"
	item.Download = magnet
	return item, nil
}

func (c *FitGirlCrawler) Crawl(page int) ([]*model.GameDownload, error) {
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: fmt.Sprintf(constant.FitGirlURL, page),
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
	updateFlags := []string{} //link+date
	doc.Find("article").Each(func(i int, s *goquery.Selection) {
		u, exist1 := s.Find(".entry-title>a").First().Attr("href")
		d, exist2 := s.Find("time").First().Attr("datetime")
		if exist1 && exist2 {
			urls = append(urls, u)
			updateFlags = append(updateFlags, fmt.Sprintf("%s%s", u, d))
		}
	})
	var res []*model.GameDownload
	for i, u := range urls {
		if db.IsFitgirlCrawled(updateFlags[i]) {
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

func (c *FitGirlCrawler) CrawlMulti(pages []int) ([]*model.GameDownload, error) {
	var res []*model.GameDownload
	for _, page := range pages {
		items, err := c.Crawl(page)
		if err != nil {
			return nil, err
		}
		res = append(res, items...)
	}
	return res, nil
}

func (c *FitGirlCrawler) CrawlAll() ([]*model.GameDownload, error) {
	var res []*model.GameDownload
	totalPageNum, err := c.GetTotalPageNum()
	if err != nil {
		return nil, err
	}
	for i := 1; i <= totalPageNum; i++ {
		items, err := c.Crawl(i)
		if err != nil {
			return nil, err
		}
		res = append(res, items...)
	}
	return res, nil
}

func (c *FitGirlCrawler) GetTotalPageNum() (int, error) {
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: fmt.Sprintf(constant.FitGirlURL, 1),
	})
	if err != nil {
		return 0, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		return 0, err
	}
	page, err := strconv.Atoi(doc.Find(".page-numbers.dots").First().Next().Text())
	if err != nil {
		return 0, err
	}
	return page, nil
}
