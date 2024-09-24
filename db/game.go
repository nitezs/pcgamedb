package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"pcgamedb/cache"
	"pcgamedb/config"
	"pcgamedb/model"
	"regexp"
	"slices"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	removeDelimiter            = regexp.MustCompile(`[:\-\+]`)
	removeRepeatingSpacesRegex = regexp.MustCompile(`\s+`)
)

func GetGameDownloadsByAuthor(regex string) ([]*model.GameDownload, error) {
	var res []*model.GameDownload
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "author", Value: primitive.Regex{Pattern: regex, Options: "i"}}}
	cursor, err := GameDownloadCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}
	if err = cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, err
}

func GetGameDownloadsByAuthorPagination(regex string, page int, pageSize int) ([]*model.GameDownload, int, error) {
	var res []*model.GameDownload
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "author", Value: primitive.Regex{Pattern: regex, Options: "i"}}}
	opts := options.Find()
	opts.SetSkip(int64((page - 1) * pageSize))
	opts.SetLimit(int64(pageSize))
	totalCount, err := GameDownloadCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	totalPage := (totalCount + int64(pageSize) - 1) / int64(pageSize)
	cursor, err := GameDownloadCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	if cursor.Err() != nil {
		return nil, 0, cursor.Err()
	}
	if err = cursor.All(ctx, &res); err != nil {
		return nil, 0, err
	}
	return res, int(totalPage), err
}

func IsGameCrawled(flag string, author string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{
		{Key: "author", Value: primitive.Regex{Pattern: author, Options: "i"}},
		{Key: "update_flag", Value: flag},
	}
	var game model.GameDownload
	err := GameDownloadCollection.FindOne(ctx, filter).Decode(&game)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return false
		}
		return false
	}
	return true
}

func IsGameCrawledByURL(url string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{
		{Key: "url", Value: url},
	}
	var game model.GameDownload
	err := GameDownloadCollection.FindOne(ctx, filter).Decode(&game)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return false
		}
		return false
	}
	return true
}

func SaveGameDownload(item *model.GameDownload) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now()
	}
	item.UpdatedAt = time.Now()
	item.Size = strings.Replace(item.Size, "gb", "GB", -1)
	item.Size = strings.Replace(item.Size, "mb", "MB", -1)
	filter := bson.M{"_id": item.ID}
	update := bson.M{"$set": item}
	opts := options.Update().SetUpsert(true)
	_, err := GameDownloadCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func SaveGameInfo(item *model.GameInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now()
	}
	item.UpdatedAt = time.Now()
	filter := bson.M{"_id": item.ID}
	update := bson.M{"$set": item}
	opts := options.Update().SetUpsert(true)
	_, err := GameInfoCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func GetAllGameDownloads() ([]*model.GameDownload, error) {
	var items []*model.GameDownload
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := GameDownloadCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var game model.GameDownload
		if err = cursor.Decode(&game); err != nil {
			return nil, err
		}
		items = append(items, &game)
	}
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}
	return items, err
}

func GetGameDownloadByUrl(url string) (*model.GameDownload, error) {
	var item model.GameDownload
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"url": url}
	err := GameDownloadCollection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return &model.GameDownload{}, nil
		}
		return nil, err
	}
	return &item, nil
}

func GetGameDownloadByID(id primitive.ObjectID) (*model.GameDownload, error) {
	var item model.GameDownload
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"_id": id}
	err := GameDownloadCollection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func GetGameDownloadsByIDs(ids []primitive.ObjectID) ([]*model.GameDownload, error) {
	var items []*model.GameDownload
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := GameDownloadCollection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var game model.GameDownload
		if err = cursor.Decode(&game); err != nil {
			return nil, err
		}
		items = append(items, &game)
	}
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}
	return items, err
}

func SearchGameInfos(name string, page int, pageSize int) ([]*model.GameInfo, int, error) {
	var items []*model.GameInfo
	name = removeDelimiter.ReplaceAllString(name, " ")
	name = removeRepeatingSpacesRegex.ReplaceAllString(name, " ")
	name = strings.TrimSpace(name)
	name = strings.Replace(name, " ", ".*", -1)
	name = fmt.Sprintf("%s.*", name)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"$or": []interface{}{
		bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: name, Options: "i"}}},
		bson.M{"aliases": bson.M{"$regex": primitive.Regex{Pattern: name, Options: "i"}}},
	}}
	totalCount, err := GameInfoCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	totalPage := (totalCount + int64(pageSize) - 1) / int64(pageSize)
	findOpts := options.Find().SetSkip(int64((page - 1) * pageSize)).SetLimit(int64(pageSize)).SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := GameInfoCollection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var game model.GameInfo
		if err = cursor.Decode(&game); err != nil {
			return nil, 0, err
		}
		game.Games, err = GetGameDownloadsByIDs(game.GameIDs)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, &game)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}
	return items, int(totalPage), nil
}

func SearchGameInfosCache(name string, page int, pageSize int) ([]*model.GameInfo, int, error) {
	type res struct {
		Items     []*model.GameInfo
		TotalPage int
	}
	name = strings.ToLower(name)
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("searchGameDetails:%s:%d:%d", name, page, pageSize)
		val, exist := cache.Get(key)
		if exist {
			var data res
			err := json.Unmarshal([]byte(val), &data)
			if err != nil {
				return nil, 0, err
			}
			return data.Items, data.TotalPage, nil
		} else {
			data, totalPage, err := SearchGameInfos(name, page, pageSize)
			if err != nil {
				return nil, 0, err
			}
			dataBytes, err := json.Marshal(res{Items: data, TotalPage: totalPage})
			if err != nil {
				return nil, 0, err
			}
			_ = cache.AddWithExpire(key, string(dataBytes), 12*time.Hour)
			return data, totalPage, nil
		}
	} else {
		return SearchGameInfos(name, page, pageSize)
	}
}

func GetGameInfoByPlatformID(platform string, id int) (*model.GameInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var filter interface{}
	switch platform {
	case "steam":
		filter = bson.M{"steam_id": id}
	case "gog":
		filter = bson.M{"gog_id": id}
	case "igdb":
		filter = bson.M{"igdb_id": id}
	}
	var game model.GameInfo
	err := GameInfoCollection.FindOne(ctx, filter).Decode(&game)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func GetUnorganizedGameDownloads(num int) ([]*model.GameDownload, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var gamesNotInDetails []*model.GameDownload
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "game_infos"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "games"},
			{Key: "as", Value: "gameDetail"},
		}}},
	}
	if num != -1 && num > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: num}})
	}
	pipeline = append(pipeline,
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "gameDetail", Value: bson.D{{Key: "$size", Value: 0}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "name", Value: 1}}}},
	)

	cursor, err := GameDownloadCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var game model.GameDownload
		if err := cursor.Decode(&game); err != nil {
			return nil, err
		}
		gamesNotInDetails = append(gamesNotInDetails, &game)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return gamesNotInDetails, nil
}

func GetGameInfoByID(id primitive.ObjectID) (*model.GameInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var game model.GameInfo
	err := GameInfoCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&game)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func DeduplicateGames() ([]primitive.ObjectID, error) {
	type queryRes struct {
		ID    string               `bson:"_id"`
		Total int                  `bson:"total"`
		IDs   []primitive.ObjectID `bson:"ids"`
	}

	var res []primitive.ObjectID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var qres []queryRes
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$download"},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "ids", Value: bson.D{{Key: "$push", Value: "$_id"}}},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "total", Value: bson.D{{Key: "$gt", Value: 1}}},
		}}},
	}
	cursor, err := GameDownloadCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &qres); err != nil {
		return nil, err
	}
	for _, item := range qres {
		idsToDelete := item.IDs[1:]
		res = append(res, idsToDelete...)
		_, err = GameDownloadCollection.DeleteMany(ctx, bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: idsToDelete}}}})
		if err != nil {
			return nil, err
		}
		cursor, err := GameInfoCollection.Find(ctx, bson.M{"games": bson.M{"$in": idsToDelete}})
		if err != nil {
			return nil, err
		}
		var infos []*model.GameInfo
		if err := cursor.All(ctx, &infos); err != nil {
			return nil, err
		}
		for _, info := range infos {
			newGames := make([]primitive.ObjectID, 0, len(info.GameIDs))
			for _, id := range info.GameIDs {
				if !slices.Contains(idsToDelete, id) {
					newGames = append(newGames, id)
				}
			}
			info.GameIDs = newGames
			if err := SaveGameInfo(info); err != nil {
				return nil, err
			}
		}
	}
	_, _ = CleanOrphanGamesInGameInfos()
	return res, nil
}

func CleanOrphanGamesInGameInfos() (map[primitive.ObjectID]primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$unwind", Value: "$games"}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "game_downloads"},
			{Key: "localField", Value: "games"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "gameDownloads"},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "gameDownloads", Value: bson.D{{Key: "$size", Value: 0}}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "game", Value: "$games"},
		}}},
	}
	cursor, err := GameInfoCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	qres := make([]struct {
		ID   primitive.ObjectID `bson:"_id"`
		Game primitive.ObjectID `bson:"game"`
	}, 0)
	if err := cursor.All(ctx, &qres); err != nil {
		return nil, err
	}
	var res = make(map[primitive.ObjectID]primitive.ObjectID)
	for _, item := range qres {
		info, err := GetGameInfoByID(item.ID)
		if err != nil {
			continue
		}
		newGames := make([]primitive.ObjectID, 0, len(info.GameIDs))
		for _, id := range info.GameIDs {
			if id != item.Game {
				newGames = append(newGames, id)
			}
		}
		info.GameIDs = newGames
		if err := SaveGameInfo(info); err != nil {
			return nil, err
		}
		res[item.ID] = item.Game
	}
	return res, nil
}

func CleanGameInfoWithEmptyGameIDs() ([]primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"games": bson.M{"$size": 0}}
	cursor, err := GameInfoCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var games []*model.GameInfo
	var res []primitive.ObjectID
	if err = cursor.All(ctx, &games); err != nil {
		return nil, err
	}
	for _, item := range games {
		res = append(res, item.ID)
	}
	_, err = GameInfoCollection.DeleteMany(ctx, filter)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetGameInfosByName(name string) ([]*model.GameInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	name = strings.TrimSpace(name)
	name = fmt.Sprintf("^%s$", name)
	filter := bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: name, Options: "i"}}}
	cursor, err := GameInfoCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var games []*model.GameInfo
	if err = cursor.All(ctx, &games); err != nil {
		return nil, err
	}
	return games, nil
}

func GetGameDownloadByRawName(name string) ([]*model.GameDownload, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	name = strings.TrimSpace(name)
	name = fmt.Sprintf("^%s$", name)
	filter := bson.M{"raw_name": bson.M{"$regex": primitive.Regex{Pattern: name, Options: "i"}}}
	cursor, err := GameDownloadCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var game []*model.GameDownload
	if err = cursor.All(ctx, &game); err != nil {
		return nil, err
	}
	return game, nil
}

func GetSameNameGameInfos() (map[string][]primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$name"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "ids", Value: bson.D{{Key: "$addToSet", Value: "$_id"}}},
		}}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "count", Value: bson.D{{Key: "$gt", Value: 1}}}}}},
	}
	cursor, err := GameInfoCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	data := make([]struct {
		Name  string               `bson:"_id"`
		Count int                  `bson:"count"`
		IDs   []primitive.ObjectID `bson:"ids"`
	}, 0)
	if err := cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	res := make(map[string][]primitive.ObjectID)
	for _, item := range data {
		res[item.Name] = item.IDs
	}
	return res, nil
}

func MergeSameNameGameInfos() error {
	games, err := GetSameNameGameInfos()
	if err != nil {
		return err
	}
	for _, ids := range games {
		var IGDBItem *model.GameInfo = nil
		otherPlatformItems := make([]*model.GameInfo, 0)
		skip := false
		for _, id := range ids {
			item, err := GetGameInfoByID(id)
			if err != nil {
				continue
			}
			if item.IGDBID != 0 {
				if IGDBItem == nil {
					IGDBItem = item
				} else {
					skip = true
					break
					// skip if there are multiple items with IGDB ID
					// not sure which item is correct
					// need deal manually
				}
			} else {
				otherPlatformItems = append(otherPlatformItems, item)
			}
		}
		if skip {
			continue
		}
		if IGDBItem != nil {
			for _, item := range otherPlatformItems {
				IGDBItem.GameIDs = append(IGDBItem.GameIDs, item.ID)
			}
			if err := SaveGameInfo(IGDBItem); err != nil {
				continue
			}
		}
	}
	return nil
}

func GetGameInfoCount() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	count, err := GameInfoCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetGameDownloadCount() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	count, err := GameDownloadCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetGameInfoWithSteamID() ([]*model.GameInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"$and": []bson.M{
		{"steam_id": bson.M{"$exists": 1}},
		{"steam_id": bson.M{"$ne": 0}},
	}}

	cursor, err := GameInfoCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var games []*model.GameInfo
	if err = cursor.All(ctx, &games); err != nil {
		return nil, err
	}
	return games, nil
}

func DeleteGameInfoByID(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := GameInfoCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func DeleteGameDownloadByID(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := GameDownloadCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	filter := bson.M{"games": bson.M{"$in": []primitive.ObjectID{id}}}
	cursor, err := GameInfoCollection.Find(ctx, filter)
	if err != nil {
		return err
	}
	var games []*model.GameInfo
	if err = cursor.All(ctx, &games); err != nil {
		return err
	}
	for _, game := range games {
		newIDs := make([]primitive.ObjectID, 0)
		for _, gameID := range game.GameIDs {
			if gameID != id {
				newIDs = append(newIDs, gameID)
			}
		}
		game.GameIDs = newIDs
		if err := SaveGameInfo(game); err != nil {
			continue
		}
	}
	return nil
}

func GetAllAuthors() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$author"},
		}}},
	}

	cursor, err := GameDownloadCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var authors []struct {
		Author string `bson:"_id"`
	}
	if err = cursor.All(ctx, &authors); err != nil {
		return nil, err
	}
	var res []string
	for _, author := range authors {
		res = append(res, author.Author)
	}
	return res, nil
}

func GetAllGameInfos() ([]*model.GameInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := GameInfoCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var res []*model.GameInfo
	if err = cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}
