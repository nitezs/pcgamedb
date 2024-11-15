package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/nitezs/pcgamedb/constant"
	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/model"
	"github.com/nitezs/pcgamedb/utils"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type ChovkaCrawler struct {
	logger *zap.Logger
}

func NewChovkaCrawler(logger *zap.Logger) *ChovkaCrawler {
	return &ChovkaCrawler{
		logger: logger,
	}
}

func (c *ChovkaCrawler) Name() string {
	return "ChovkaCrawler"
}

func (c *ChovkaCrawler) CrawlByUrl(url string) (*model.GameDownload, error) {
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
	item.Name = ChovkaFormatter(item.RawName)
	item.Author = "Chovka"
	item.UpdateFlag = item.RawName
	downloadURL := doc.Find(".download-torrent").AttrOr("href", "")
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

func (c *ChovkaCrawler) Crawl(page int) ([]*model.GameDownload, error) {
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: fmt.Sprintf(constant.RepackInfoURL, page),
	})
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		return nil, err
	}
	urls := []string{}
	updateFlags := []string{}
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
		if db.IsChovkaCrawled(updateFlags[i]) {
			continue
		}
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

func (c *ChovkaCrawler) CrawlMulti(pages []int) ([]*model.GameDownload, error) {
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

func (c *ChovkaCrawler) CrawlAll() ([]*model.GameDownload, error) {
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

func (c *ChovkaCrawler) GetTotalPageNum() (int, error) {
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: fmt.Sprintf(constant.RepackInfoURL, 1),
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

func ChovkaFormatter(name string) string {
	idx := strings.Index(name, "| RePack")
	if idx != -1 {
		name = name[:idx]
	}
	idx = strings.Index(name, "| GOG")
	if idx != -1 {
		name = name[:idx]
	}
	idx = strings.Index(name, "| Portable")
	if idx != -1 {
		name = name[:idx]
	}
	return strings.TrimSpace(name)
}
