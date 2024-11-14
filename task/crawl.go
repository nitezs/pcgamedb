package task

import (
	"net/http"
	"net/url"
	"pcgamedb/config"
	"pcgamedb/crawler"
	"pcgamedb/model"
	"pcgamedb/utils"

	"go.uber.org/zap"
)

func Crawl(logger *zap.Logger) {
	var games []*model.GameDownload
	var crawlerMap = crawler.BuildCrawlerMap(logger)
	for _, item := range crawlerMap {
		if c, ok := item.(crawler.PagedCrawler); ok {
			g, err := c.CrawlMulti([]int{1, 2, 3})
			if err != nil {
				logger.Error("Failed to crawl games", zap.Error(err))
			}
			games = append(games, g...)
		} else if c, ok := item.(crawler.SimpleCrawler); ok {
			g, err := c.CrawlAll()
			if err != nil {
				logger.Error("Failed to crawl games", zap.Error(err))
			}
			games = append(games, g...)
		}
	}
	logger.Info("Crawled finished", zap.Int("count", len(games)))
	for _, game := range games {
		logger.Info(
			"Crawled game",
			zap.String("name", game.RawName),
			zap.String("author", game.Author),
			zap.String("url", game.Url),
		)
	}
	Clean(logger)
	for _, u := range config.Config.Webhooks.CrawlTask {
		_, err := url.Parse(u)
		if err != nil {
			logger.Error("Invalid webhook url", zap.String("url", u), zap.Error(err))
			continue
		}
		logger.Info("webhook triggered", zap.String("task", "crawl"), zap.String("url", u))
		_, err = utils.Fetch(utils.FetchConfig{
			Url:    u,
			Method: http.MethodPost,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Data: games,
		})
		if err != nil {
			logger.Error("Failed to trigger webhook", zap.String("task", "crawl"), zap.String("url", u), zap.Error(err))
		}
	}
}
