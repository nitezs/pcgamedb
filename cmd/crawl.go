package cmd

import (
	"errors"
	"fmt"
	"pcgamedb/crawler"
	"pcgamedb/log"
	"pcgamedb/utils"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Long:  "Crawl games from specific platforms",
	Short: "Crawl games from specific platforms",
	Run:   crawlRun,
}

type CrawlCommandConfig struct {
	Source string
	Page   string
	All    bool
	Num    int
}

var crawlCmdCfg CrawlCommandConfig

var crawlerMap = map[string]crawler.Crawler{}

func init() {
	crawlerMap = crawler.BuildCrawlerMap(log.Logger)
	allCrawlerBuilder := strings.Builder{}
	paginationCrwalerBuilder := strings.Builder{}
	noPaginationCrawlerBuilder := strings.Builder{}
	for k, v := range crawlerMap {
		allCrawlerBuilder.WriteString(k)
		allCrawlerBuilder.WriteString(",")
		if _, ok := v.(crawler.PagedCrawler); ok {
			paginationCrwalerBuilder.WriteString(k)
			paginationCrwalerBuilder.WriteString(",")
		} else if _, ok := v.(crawler.SimpleCrawler); ok {
			noPaginationCrawlerBuilder.WriteString(k)
			noPaginationCrawlerBuilder.WriteString(",")
		}
	}
	crawlCmd.Flags().StringVarP(&crawlCmdCfg.Source, "source", "s", "", fmt.Sprintf("source to crawl (%s)", strings.Trim(allCrawlerBuilder.String(), ",")))
	crawlCmd.Flags().StringVarP(&crawlCmdCfg.Page, "pages", "p", "1", fmt.Sprintf("pages to crawl (1,2,3 or 1-3) (%s)", strings.Trim(paginationCrwalerBuilder.String(), ",")))
	crawlCmd.Flags().BoolVarP(&crawlCmdCfg.All, "all", "a", false, "crawl all page")
	crawlCmd.Flags().IntVarP(&crawlCmdCfg.Num, "num", "n", -1, fmt.Sprintf("number of items to process (%s)", strings.Trim(noPaginationCrawlerBuilder.String(), ",")))
	RootCmd.AddCommand(crawlCmd)
}

func crawlRun(cmd *cobra.Command, args []string) {
	crawlCmdCfg.Source = strings.ToLower(crawlCmdCfg.Source)

	if crawlCmdCfg.Source == "" {
		log.Logger.Error("Source is required")
		return
	}

	item, ok := crawlerMap[crawlCmdCfg.Source]
	if !ok {
		log.Logger.Error("Invalid source", zap.String("source", crawlCmdCfg.Source))
		return
	}

	if c, ok := item.(crawler.PagedCrawler); ok {
		if crawlCmdCfg.All {
			_, err := c.CrawlAll()
			if err != nil {
				log.Logger.Error("Crawl error", zap.Error(err))
				return
			}
		} else {
			pages, err := pagination(crawlCmdCfg.Page)
			if err != nil {
				log.Logger.Error("Invalid page", zap.String("page", crawlCmdCfg.Page))
				return
			}
			_, err = c.CrawlMulti(pages)
			if err != nil {
				log.Logger.Error("Crawl error", zap.Error(err))
				return
			}
		}
	} else if c, ok := item.(crawler.SimpleCrawler); ok {
		if crawlCmdCfg.All {
			_, err := c.CrawlAll()
			if err != nil {
				log.Logger.Error("Crawl error", zap.Error(err))
				return
			}
		} else {
			_, err := c.Crawl(crawlCmdCfg.Num)
			if err != nil {
				log.Logger.Error("Crawl error", zap.Error(err))
				return
			}
		}
	}
}

func pagination(pageStr string) ([]int, error) {
	if pageStr == "" {
		return nil, errors.New("empty page")
	}
	var pages []int
	pageSlice := strings.Split(pageStr, ",")
	for i := 0; i < len(pageSlice); i++ {
		if strings.Contains(pageSlice[i], "-") {
			pageRange := strings.Split(pageSlice[i], "-")
			start, err := strconv.Atoi(pageRange[0])
			if err != nil {
				return nil, err
			}
			end, err := strconv.Atoi(pageRange[1])
			if err != nil {
				return nil, err
			}
			if start > end {
				return nil, err
			}
			for j := start; j <= end; j++ {
				pages = append(pages, j)
			}
		} else {
			p, err := strconv.Atoi(pageSlice[i])
			if err != nil {
				log.Logger.Error("Invalid page", zap.String("page", pageSlice[i]))
				return nil, err
			}
			pages = append(pages, p)
		}
	}
	return utils.Unique(pages), nil
}
