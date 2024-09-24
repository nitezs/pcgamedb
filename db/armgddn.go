package db

import "pcgamedb/model"

func IsARMGDDNCrawled(flag string) bool {
	return IsGameCrawled(flag, "armgddn")
}

func GetARMGDDNGameDownloads() ([]*model.GameDownload, error) {
	return GetGameDownloadsByAuthor("armgddn")
}
