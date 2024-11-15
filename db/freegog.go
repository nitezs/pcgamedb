package db

import (
	"github.com/nitezs/pcgamedb/model"
)

func GetFreeGOGGameDownloads() ([]*model.GameDownload, error) {
	return GetGameDownloadsByAuthor("freegog")
}
func IsFreeGOGCrawled(flag string) bool {
	return IsGameCrawled(flag, "freegog")
}
