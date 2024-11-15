package db

import "github.com/nitezs/pcgamedb/model"

func GetFitgirlAllGameDownloads() ([]*model.GameDownload, error) {
	return GetGameDownloadsByAuthor("fitgirl")
}

func IsFitgirlCrawled(flag string) bool {
	return IsGameCrawled(flag, "fitgirl")
}
