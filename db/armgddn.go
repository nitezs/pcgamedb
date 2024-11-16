package db

import "github.com/nitezs/pcgamedb/model"

func IsARMGDDNCrawled(flag string) bool {
	return IsGameCrawled(flag, "armgddn")
}

func GetARMGDDNGameItems() ([]*model.GameItem, error) {
	return GetGameItemsByAuthor("armgddn")
}
