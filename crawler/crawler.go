package crawler

import (
	"pcgamedb/model"

	"go.uber.org/zap"
)

type Crawler interface {
	Crawl(int) ([]*model.GameDownload, error)
	CrawlAll() ([]*model.GameDownload, error)
}

type SimpleCrawler interface {
	Crawler
}

type PagedCrawler interface {
	Crawler
	CrawlMulti([]int) ([]*model.GameDownload, error)
	GetTotalPageNum() (int, error)
}

func BuildCrawlerMap(logger *zap.Logger) map[string]Crawler {
	return map[string]Crawler{
		"fitgirl":  NewFitGirlCrawler(logger),
		"dodi":     NewDODICrawler(logger),
		"kaoskrew": NewKaOsKrewCrawler(logger),
		// "freegog":   NewFreeGOGCrawler(logger),
		"xatab":     NewXatabCrawler(logger),
		"onlinefix": NewOnlineFixCrawler(logger),
		"steamrip":  NewSteamRIPCrawler(logger),
		// "armgddn":   NewARMGDDNCrawler(logger),
		"goggames": NewGOGGamesCrawler(logger),
		"chovka":   NewChovkaCrawler(logger),
		// "gnarly":   NewGnarlyCrawler(logger),
	}
}
