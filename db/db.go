package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nitezs/pcgamedb/config"
	"github.com/nitezs/pcgamedb/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	gameDownloadCollectionName = "games"
	gameInfoCollectionName     = "game_infos"
)

var (
	mongoDB            *mongo.Client
	mutx               = &sync.RWMutex{}
	GameItemCollection = &CustomCollection{
		collName: gameDownloadCollectionName,
	}
	GameInfoCollection = &CustomCollection{
		collName: gameInfoCollectionName,
	}
)

func connect() {
	if !config.Config.DatabaseAvaliable {
		log.Logger.Panic("Missing database configuration information")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(fmt.Sprintf(
		"mongodb://%s:%s@%s:%v",
		config.Config.Database.User,
		config.Config.Database.Password,
		config.Config.Database.Host,
		config.Config.Database.Port,
	))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Logger.Panic("Failed to connect to MongoDB", zap.Error(err))
	}
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Logger.Panic("Failed to ping MongoDB", zap.Error(err))
	}
	log.Logger.Info("Connected to MongoDB")
	mongoDB = client

	gameDownloadCollection := mongoDB.Database(config.Config.Database.Database).Collection(gameDownloadCollectionName)
	gameInfoCollection := mongoDB.Database(config.Config.Database.Database).Collection(gameInfoCollectionName)

	nameIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
		},
	}
	authorIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "author", Value: 1},
		},
	}
	gamesIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "games", Value: 1},
		},
	}
	searchIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: "text"}, {Key: "aliases", Value: "text"}},
	}
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = gameDownloadCollection.Indexes().CreateOne(ctx, nameIndex)
	if err != nil {
		log.Logger.Error("Failed to create index", zap.Error(err))
	}
	_, err = gameDownloadCollection.Indexes().CreateOne(ctx, authorIndex)
	if err != nil {
		log.Logger.Error("Failed to create index", zap.Error(err))
	}
	_, err = gameInfoCollection.Indexes().CreateOne(ctx, gamesIndex)
	if err != nil {
		log.Logger.Error("Failed to create index", zap.Error(err))
	}
	_, err = gameInfoCollection.Indexes().CreateOne(ctx, searchIndex)
	if err != nil {
		log.Logger.Error("Failed to create index", zap.Error(err))
	}
}

func CheckConnect() {
	mutx.RLock()
	if mongoDB != nil {
		mutx.RUnlock()
		return
	}
	mutx.RUnlock()

	mutx.Lock()
	if mongoDB == nil {
		connect()
	}
	mutx.Unlock()
}

func HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mongoDB.Ping(ctx, nil)
}
