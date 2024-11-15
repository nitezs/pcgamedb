package db

import (
	"github.com/nitezs/pcgamedb/model"
)

func GetXatabGameDownloads() ([]*model.GameDownload, error) {
	return GetGameDownloadsByAuthor("xatab")
}

func IsXatabCrawled(flag string) bool {
	return IsGameCrawled(flag, "xatab")
}
