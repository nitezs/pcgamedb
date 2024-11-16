package db

import (
	"github.com/nitezs/pcgamedb/model"
)

func GetFreeGOGGameItems() ([]*model.GameItem, error) {
	return GetGameItemsByAuthor("freegog")
}
func IsFreeGOGCrawled(flag string) bool {
	return IsGameCrawled(flag, "freegog")
}
