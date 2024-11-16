package db

import (
	"github.com/nitezs/pcgamedb/model"
)

func GetXatabGameItems() ([]*model.GameItem, error) {
	return GetGameItemsByAuthor("xatab")
}

func IsXatabCrawled(flag string) bool {
	return IsGameCrawled(flag, "xatab")
}
