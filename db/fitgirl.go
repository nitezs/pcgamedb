package db

import "github.com/nitezs/pcgamedb/model"

func GetFitgirlAllGameItems() ([]*model.GameItem, error) {
	return GetGameItemsByAuthor("fitgirl")
}

func IsFitgirlCrawled(flag string) bool {
	return IsGameCrawled(flag, "fitgirl")
}
