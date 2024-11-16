package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nitezs/pcgamedb/model"
	"go.mongodb.org/mongo-driver/bson"
)

func Export() ([]byte, []byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	infos := []model.GameInfo{}
	games := []model.GameItem{}
	cursor, err := GameInfoCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &infos); err != nil {
		return nil, nil, err
	}
	cursor, err = GameItemCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &games); err != nil {
		return nil, nil, err
	}
	infoJson, err := json.Marshal(infos)
	if err != nil {
		return nil, nil, err
	}
	gameJson, err := json.Marshal(games)
	if err != nil {
		return nil, nil, err
	}
	return infoJson, gameJson, nil
}
