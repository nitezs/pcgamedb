package db

import (
	"pcgamedb/model"
)

func GetFreeGOGGameDownloads() ([]*model.GameDownload, error) {
	return GetGameDownloadsByAuthor("freegog")
}
func IsFreeGOGCrawled(flag string) bool {
	return IsGameCrawled(flag, "freegog")
}
