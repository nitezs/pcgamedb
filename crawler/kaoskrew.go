package crawler

import (
	"pcgamedb/model"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

const KaOsKrewName string = "KaOsKrew-torrents"

type KaOsKrewCrawler struct {
	logger  *zap.Logger
	crawler s1337xCrawler
}

func NewKaOsKrewCrawler(logger *zap.Logger) *KaOsKrewCrawler {
	return &KaOsKrewCrawler{
		logger: logger,
		crawler: *New1337xCrawler(
			KaOsKrewName,
			KaOsKrewFormatter,
			logger,
		),
	}
}

func (c *KaOsKrewCrawler) Name() string {
	return "KaOsKrewCrawler"
}

func (c *KaOsKrewCrawler) Crawl(page int) ([]*model.GameDownload, error) {
	return c.crawler.Crawl(page)
}

func (c *KaOsKrewCrawler) CrawlByUrl(url string) (*model.GameDownload, error) {
	return c.crawler.CrawlByUrl(url)
}

func (c *KaOsKrewCrawler) CrawlMulti(pages []int) ([]*model.GameDownload, error) {
	return c.crawler.CrawlMulti(pages)
}

func (c *KaOsKrewCrawler) CrawlAll() ([]*model.GameDownload, error) {
	return c.crawler.CrawlAll()
}

func (c *KaOsKrewCrawler) GetTotalPageNum() (int, error) {
	return c.crawler.GetTotalPageNum()
}

var kaOsKrewRegexps = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\.REPACK2?-KaOs`),
	regexp.MustCompile(`(?i)\.UPDATE-KaOs`),
	regexp.MustCompile(`(?i)v\.?\d+(\.\d+)*|Build\.\d+`),
	regexp.MustCompile(`(?i)\.MULTi\d+`),
	regexp.MustCompile(`(?i)\sgoty`),
}

func KaOsKrewFormatter(name string) string {
	if index := kaOsKrewRegexps[2].FindIndex([]byte(name)); index != nil {
		name = name[:index[0]]
	}
	for _, re := range kaOsKrewRegexps {
		name = re.ReplaceAllString(name, "")
	}
	name = strings.Replace(name, ".", " ", -1)
	return strings.TrimSpace(name)
}
