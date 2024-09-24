package crawler

import (
	"bytes"
	"pcgamedb/constant"
	"pcgamedb/db"
	"pcgamedb/model"
	"pcgamedb/utils"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type GnarlyCrawler struct {
	logger *zap.Logger
}

func NewGnarlyCrawler(logger *zap.Logger) *GnarlyCrawler {
	return &GnarlyCrawler{
		logger: logger,
	}
}

func (c *GnarlyCrawler) Crawl(num int) ([]*model.GameDownload, error) {
	var res []*model.GameDownload
	count := 0
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.GnarlyURL,
	})
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		return nil, err
	}
	sizeRegex := regexp.MustCompile(`\[(\d+)\s(GB|MB)\]`)
	pElementHtml := make([]string, 0)
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		pElementHtml = append(pElementHtml, s.Text())
	})
	for _, s := range pElementHtml {
		if strings.Contains(s, "https://bin.0xfc.de/") {
			lines := strings.Split(s, "\n")
			for i := 0; i < len(lines); i++ {
				if strings.Contains(lines[i], "[Gnarly Repacks]") {
					i++
					if strings.Contains(lines[i], "https://bin.0xfc.de/") {
						if count == num {
							return res, nil
						}
						if db.IsGnarlyCrawled(lines[i-1]) {
							continue
						}
						item, err := db.GetGameDownloadByUrl(lines[i])
						if err != nil {
							continue
						}
						sizeRegexRes := sizeRegex.FindStringSubmatch(lines[i])
						if len(sizeRegexRes) == 3 {
							item.Size = sizeRegexRes[1] + " " + sizeRegexRes[2]
						}
						c.logger.Info("Crawling", zap.String("Name", lines[i-1]))
						item.RawName = lines[i-1]
						item.Url = constant.GnarlyURL
						item.Author = "Gnarly"
						item.Name = GnarlyFormatter(item.RawName)
						download, err := utils.DecryptPrivateBin(lines[i], "gnarly")
						if err != nil {
							continue
						}
						item.Download = download
						item.UpdateFlag = item.RawName
						res = append(res, item)
						count++
						info, err := OrganizeGameDownload(item)
						if err != nil {
							continue
						}
						err = db.SaveGameInfo(info)
						if err != nil {
							c.logger.Warn("Failed to save game info", zap.Error(err))
							continue
						}
					}
				}
			}
		}
	}
	return res, nil
}

func (c *GnarlyCrawler) CrawlAll() ([]*model.GameDownload, error) {
	return c.Crawl(-1)
}

var parenthesesRegex = regexp.MustCompile(`\(([^)]+)\)`)

func GnarlyFormatter(name string) string {
	name = name[:strings.Index(name, " [Gnarly Repacks]")]
	name = parenthesesRegex.ReplaceAllString(name, "")
	return strings.TrimSpace(name)
}
