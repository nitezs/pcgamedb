package db

import (
	"github.com/nitezs/pcgamedb/model"
)

func GetOnlineFixGameItems() ([]*model.GameItem, error) {
	return GetGameItemsByAuthor("onlinefix")
}

func IsOnlineFixCrawled(flag string) bool {
	return IsGameCrawled(flag, "onlinefix")
}
