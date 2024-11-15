package db

import (
	"github.com/nitezs/pcgamedb/model"
)

func GetDODIGameDownloads() ([]*model.GameDownload, error) {
	return GetGameDownloadsByAuthor("dodi")
}

func GetKaOsKrewGameDownloads() ([]*model.GameDownload, error) {
	return GetGameDownloadsByAuthor("kaoskrew")
}
