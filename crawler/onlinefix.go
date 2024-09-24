package crawler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"pcgamedb/config"
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

type OnlineFixCrawler struct {
	logger  *zap.Logger
	cookies map[string]string
}

func NewOnlineFixCrawler(logger *zap.Logger) *OnlineFixCrawler {
	return &OnlineFixCrawler{
		logger:  logger,
		cookies: map[string]string{},
	}
}

func (c *OnlineFixCrawler) Crawl(page int) ([]*model.GameDownload, error) {
	if !config.Config.OnlineFixAvaliable {
		c.logger.Error("Need Online Fix account")
		return nil, errors.New("Online Fix is not available")
	}
	if len(c.cookies) == 0 {
		err := c.login()
		if err != nil {
			c.logger.Error("Failed to login", zap.Error(err))
			return nil, err
		}
	}
	requestURL := fmt.Sprintf("%s/page/%d/", constant.OnlineFixURL, page)
	resp, err := utils.Fetch(utils.FetchConfig{
		Url:     requestURL,
		Cookies: c.cookies,
		Headers: map[string]string{
			"Referer": constant.OnlineFixURL,
		},
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
	doc.Find("article.news").Each(func(i int, s *goquery.Selection) {
		urls = append(urls, s.Find(".big-link").First().AttrOr("href", ""))
		updateFlags = append(
			updateFlags,
			s.Find(".big-link").First().AttrOr("href", "")+
				s.Find("time").Text(),
		)
	})

	var res []*model.GameDownload
	for i, u := range urls {
		if db.IsOnlineFixCrawled(updateFlags[i]) {
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

func (c *OnlineFixCrawler) CrawlByUrl(url string) (*model.GameDownload, error) {
	if len(c.cookies) == 0 {
		err := c.login()
		if err != nil {
			c.logger.Error("Failed to login", zap.Error(err))
			return nil, err
		}
	}
	resp, err := utils.Fetch(utils.FetchConfig{
		Url:     url,
		Cookies: c.cookies,
		Headers: map[string]string{
			"Referer": constant.OnlineFixURL,
		},
	})
	if err != nil {
		return nil, err
	}
	titleRegex := regexp.MustCompile(`(?i)<h1.*?>(.*?)</h1>`)
	titleRegexRes := titleRegex.FindAllStringSubmatch(string(resp.Data), -1)
	if len(titleRegexRes) == 0 {
		return nil, errors.New("Failed to find title")
	}
	downloadRegex := regexp.MustCompile(`(?i)<a[^>]*\bhref="([^"]+)"[^>]*>(Скачать Torrent|Скачать торрент)</a>`)
	downloadRegexRes := downloadRegex.FindAllStringSubmatch(string(resp.Data), -1)
	if len(downloadRegexRes) == 0 {
		return nil, errors.New("Failed to find download button")
	}
	item, err := db.GetGameDownloadByUrl(url)
	if err != nil {
		return nil, err
	}
	item.RawName = titleRegexRes[0][1]
	item.Name = OnlineFixFormatter(item.RawName)
	item.Url = url
	item.Author = "OnlineFix"
	item.Size = "0"
	resp, err = utils.Fetch(utils.FetchConfig{
		Url:     downloadRegexRes[0][1],
		Cookies: c.cookies,
		Headers: map[string]string{
			"Referer": url,
		},
	})
	if err != nil {
		return nil, err
	}
	if strings.Contains(downloadRegexRes[0][1], "uploads.online-fix.me") {
		magnetRegex := regexp.MustCompile(`(?i)"(.*?).torrent"`)
		magnetRegexRes := magnetRegex.FindAllStringSubmatch(string(resp.Data), -1)
		if len(magnetRegexRes) == 0 {
			return nil, errors.New("Failed to find magnet")
		}
		resp, err = utils.Fetch(utils.FetchConfig{
			Url:     downloadRegexRes[0][1] + strings.Trim(magnetRegexRes[0][0], "\""),
			Cookies: c.cookies,
			Headers: map[string]string{
				"Referer": url,
			},
		})
		if err != nil {
			return nil, err
		}
		item.Download, item.Size, err = utils.ConvertTorrentToMagnet(resp.Data)
		if err != nil {
			return nil, err
		}
	} else if strings.Contains(downloadRegexRes[0][1], "online-fix.me/ext") {
		if strings.Contains(string(resp.Data), "mega.nz") {
			if !config.Config.MegaAvaliable {
				return nil, errors.New("Mega is not avaliable")
			}
			megaRegex := regexp.MustCompile(`(?i)location.href=\\'([^\\']*)\\'`)
			megaRegexRes := megaRegex.FindAllStringSubmatch(string(resp.Data), -1)
			if len(megaRegexRes) == 0 {
				return nil, errors.New("Failed to find download link")
			}
			path, files, err := utils.MegaDownload(megaRegexRes[0][1], "torrent")
			if err != nil {
				return nil, err
			}
			torrent := ""
			for _, file := range files {
				if strings.HasSuffix(file, ".torrent") {
					torrent = file
					break
				}
			}
			dataBytes, err := os.ReadFile(torrent)
			if err != nil {
				return nil, err
			}
			item.Download, item.Size, err = utils.ConvertTorrentToMagnet(dataBytes)
			if err != nil {
				return nil, err
			}
			_ = os.RemoveAll(path)
		} else {
			return nil, errors.New("Failed to find download link")
		}
	} else {
		return nil, errors.New("Failed to find download link")
	}
	return item, nil
}

func (c *OnlineFixCrawler) CrawlMulti(pages []int) ([]*model.GameDownload, error) {
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

func (c *OnlineFixCrawler) CrawlAll() ([]*model.GameDownload, error) {
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

func (c *OnlineFixCrawler) GetTotalPageNum() (int, error) {
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.OnlineFixURL,
		Headers: map[string]string{
			"Referer": constant.OnlineFixURL,
		},
	})
	if err != nil {
		return 0, err
	}
	pageRegex := regexp.MustCompile(`(?i)<a href="https://online-fix.me/page/(\d+)/">.*?</a>`)
	pageRegexRes := pageRegex.FindAllStringSubmatch(string(resp.Data), -1)
	if len(pageRegexRes) == 0 {
		return 0, err
	}
	totalPageNum, err := strconv.Atoi(pageRegexRes[len(pageRegexRes)-2][1])
	if err != nil {
		return 0, err
	}
	return totalPageNum, nil
}

type csrf struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

func (c *OnlineFixCrawler) login() error {
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.OnlineFixCSRFURL,
		Headers: map[string]string{
			"X-Requested-With": "XMLHttpRequest",
			"Referer":          constant.OnlineFixURL,
		},
	})
	if err != nil {
		return err
	}
	var csrf csrf
	if err = json.Unmarshal(resp.Data, &csrf); err != nil {
		return err
	}

	for _, cookie := range resp.Cookie {
		c.cookies[cookie.Name] = cookie.Value
	}
	params := url.Values{}
	params.Add("login_name", config.Config.OnlineFix.User)
	params.Add("login_password", config.Config.OnlineFix.Password)
	params.Add(csrf.Field, csrf.Value)
	params.Add("login", "submit")
	resp, err = utils.Fetch(utils.FetchConfig{
		Url:     constant.OnlineFixURL,
		Method:  "POST",
		Cookies: c.cookies,
		Headers: map[string]string{
			"Origin":       constant.OnlineFixURL,
			"Content-Type": "application/x-www-form-urlencoded",
			"Referer":      constant.OnlineFixURL,
		},
		Data: params,
	})
	if err != nil {
		return err
	}
	for _, cookie := range resp.Cookie {
		c.cookies[cookie.Name] = cookie.Value
	}
	return nil
}

func OnlineFixFormatter(name string) string {
	name = strings.Replace(name, "по сети", "", -1)
	reg1 := regexp.MustCompile(`(?i)\(.*?\)`)
	name = reg1.ReplaceAllString(name, "")
	return strings.TrimSpace(name)
}
