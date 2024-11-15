package crawler

import (
	"regexp"
	"strings"

	"github.com/nitezs/pcgamedb/model"
	"github.com/nitezs/pcgamedb/utils"

	"go.uber.org/zap"
)

const DODIName string = "DODI-torrents"

type DODICrawler struct {
	logger  *zap.Logger
	crawler s1337xCrawler
}

func NewDODICrawler(logger *zap.Logger) *DODICrawler {
	return &DODICrawler{
		logger: logger,
		crawler: *New1337xCrawler(
			DODIName,
			DODIFormatter,
			logger,
		),
	}
}

func (c *DODICrawler) Name() string {
	return "DODICrawler"
}

func (c *DODICrawler) Crawl(page int) ([]*model.GameDownload, error) {
	return c.crawler.Crawl(page)
}

func (c *DODICrawler) CrawlByUrl(url string) (*model.GameDownload, error) {
	return c.crawler.CrawlByUrl(url)
}

func (c *DODICrawler) CrawlMulti(pages []int) ([]*model.GameDownload, error) {
	return c.crawler.CrawlMulti(pages)
}

func (c *DODICrawler) CrawlAll() ([]*model.GameDownload, error) {
	return c.crawler.CrawlAll()
}

func (c *DODICrawler) GetTotalPageNum() (int, error) {
	return c.crawler.GetTotalPageNum()
}

var dodiRegexps = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\s{2,}`),
	regexp.MustCompile(`(?i)[\-\+]\s?[^:\-]*?\s(Edition|Bundle|Pack|Set|Remake|Collection)`),
}

func DODIFormatter(name string) string {
	name = strings.Replace(name, "- [DODI Repack]", "", -1)
	name = strings.Replace(name, "- Campaign Remastered", "", -1)
	name = strings.Replace(name, "- Remastered", "", -1)
	if index := strings.Index(name, "+"); index != -1 {
		name = name[:index]
	}
	if index := strings.Index(name, "–"); index != -1 {
		name = name[:index]
	}
	if index := strings.Index(name, "("); index != -1 {
		name = name[:index]
	}
	if index := strings.Index(name, "["); index != -1 {
		name = name[:index]
	}
	if index := strings.Index(name, "- AiO"); index != -1 {
		name = name[:index]
	}
	if index := strings.Index(name, "- All In One"); index != -1 {
		name = name[:index]
	}
	for _, re := range dodiRegexps {
		name = strings.TrimSpace(re.ReplaceAllString(name, ""))
	}
	name = strings.TrimSpace(name)
	name = strings.Replace(name, "- Portable", "", -1)
	name = strings.Replace(name, "- Remastered", "", -1)

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
