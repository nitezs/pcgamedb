package db

import "github.com/nitezs/pcgamedb/model"

func IsARMGDDNCrawled(flag string) bool {
	return IsGameCrawled(flag, "armgddn")
}

func GetARMGDDNGameDownloads() ([]*model.GameDownload, error) {
	return GetGameDownloadsByAuthor("armgddn")
}
