package db

import (
	"pcgamedb/model"
)

func GetXatabGameDownloads() ([]*model.GameDownload, error) {
	return GetGameDownloadsByAuthor("xatab")
}

func IsXatabCrawled(flag string) bool {
	return IsGameCrawled(flag, "xatab")
}
