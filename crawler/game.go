package crawler

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/model"
	"github.com/nitezs/pcgamedb/utils"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GenerateGameInfo(platform string, id int) (*model.GameInfo, error) {
	switch platform {
	case "steam":
		return GenerateSteamGameInfo(id)
	case "igdb":
		return GenerateIGDBGameInfo(id)
	default:
		return nil, errors.New("Invalid ID type")
	}
}

func OrganizeGameItem(game *model.GameItem) (*model.GameInfo, error) {
	item, err := OrganizeGameItemWithIGDB(0, game)
	if err == nil {
		if item.SteamID == 0 {
			// get steam id from igdb
			steamID, err := GetSteamIDByIGDBIDCache(item.IGDBID)
			if err == nil {
				item.SteamID = steamID
			}
			return item, nil
		}
	}
	item, err = OrganizeGameItemWithSteam(0, game)
	if err == nil {
		if item.IGDBID == 0 {
			igdbID, err := GetIGDBIDBySteamIDCache(item.SteamID)
			if err == nil {
				item.IGDBID = igdbID
			}
		}
		return item, nil
	}
	return nil, err
}

func AddGameInfoManually(gameID primitive.ObjectID, platform string, plateformID int) (*model.GameInfo, error) {
	info, err := GenerateGameInfo(platform, plateformID)
	if err != nil {
		return nil, err
	}
	info.GameIDs = append(info.GameIDs, gameID)
	info.GameIDs = utils.Unique(info.GameIDs)
	return info, db.SaveGameInfo(info)
}

func OrganizeGameItemManually(gameID primitive.ObjectID, platform string, platformID int) (*model.GameInfo, error) {
	info, err := db.GetGameInfoByPlatformID(platform, platformID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			info, err = AddGameInfoManually(gameID, platform, platformID)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	info.GameIDs = append(info.GameIDs, gameID)
	info.GameIDs = utils.Unique(info.GameIDs)
	err = db.SaveGameInfo(info)
	if err != nil {
		return nil, err
	}
	if platform == "igdb" {
		steamID, err := GetSteamIDByIGDBIDCache(platformID)
		if err == nil {
			info.SteamID = steamID
		}
	}
	if platform == "steam" {
		igdbID, err := GetIGDBIDBySteamIDCache(platformID)
		if err == nil {
			info.IGDBID = igdbID
		}
	}
	return info, nil
}

func FormatName(name string) string {
	name = regexp.MustCompile(`(?i)[\w’'-]+\s(Edition|Vision|Collection|Bundle|Pack|Deluxe)`).ReplaceAllString(name, " ")
	name = regexp.MustCompile(`(?i)GOTY`).ReplaceAllString(name, "")
	name = regexp.MustCompile(`(?i)nsw for pc`).ReplaceAllString(name, "")
	name = regexp.MustCompile(`\([^\)]+\)`).ReplaceAllString(name, "")
	name = regexp.MustCompile(`\s+`).ReplaceAllString(name, " ")
	name = strings.Replace(name, ": Remastered", "", -1)
	name = strings.Replace(name, ": Remaster", "", -1)
	name = strings.TrimSpace(name)
	name = strings.Trim(name, ":")
	return name
}

func SupplementPlatformIDToGameInfo(logger *zap.Logger) error {
	infos, err := db.GetAllGameInfos()
	if err != nil {
		return err
	}
	for _, info := range infos {
		changed := false
		if info.IGDBID != 0 && info.SteamID == 0 {
			steamID, err := GetSteamIDByIGDBIDCache(info.IGDBID)
			time.Sleep(time.Millisecond * 100)
			if err != nil {
				continue
			}
			info.SteamID = steamID
			changed = true
		}
		if info.SteamID != 0 && info.IGDBID == 0 {
			igdbID, err := GetIGDBIDBySteamIDCache(info.SteamID)
			time.Sleep(time.Millisecond * 100)
			if err != nil {
				continue
			}
			info.IGDBID = igdbID
			changed = true
		}
		if changed {
			logger.Info("Supplemented platform id for game info", zap.String("name", info.Name), zap.Int("igdb", int(info.IGDBID)), zap.Int("steam", int(info.SteamID)))
			_ = db.SaveGameInfo(info)
		}
	}
	return nil
}
