package db

import (
	"github.com/nitezs/pcgamedb/model"
)

func GetDODIGameItems() ([]*model.GameItem, error) {
	return GetGameItemsByAuthor("dodi")
}

func GetKaOsKrewGameItems() ([]*model.GameItem, error) {
	return GetGameItemsByAuthor("kaoskrew")
}
