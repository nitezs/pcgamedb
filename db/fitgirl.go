package db

import "pcgamedb/model"

func GetFitgirlAllGameDownloads() ([]*model.GameDownload, error) {
	return GetGameDownloadsByAuthor("fitgirl")
}

func IsFitgirlCrawled(flag string) bool {
	return IsGameCrawled(flag, "fitgirl")
}
