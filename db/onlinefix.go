package db

import (
	"pcgamedb/model"
)

func GetOnlineFixGameDownloads() ([]*model.GameDownload, error) {
	return GetGameDownloadsByAuthor("onlinefix")
}

func IsOnlineFixCrawled(flag string) bool {
	return IsGameCrawled(flag, "onlinefix")
}
