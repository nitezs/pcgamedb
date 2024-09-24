package db

import (
	"pcgamedb/model"
)

func GetDODIGameDownloads() ([]*model.GameDownload, error) {
	return GetGameDownloadsByAuthor("dodi")
}

func GetKaOsKrewGameDownloads() ([]*model.GameDownload, error) {
	return GetGameDownloadsByAuthor("kaoskrew")
}
