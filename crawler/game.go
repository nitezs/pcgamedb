package crawler

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/model"
	"github.com/nitezs/pcgamedb/utils"

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

func OrganizeGameDownload(game *model.GameDownload) (*model.GameInfo, error) {
	item, err := OrganizeGameDownloadWithIGDB(0, game)
	if err == nil {
		if item.SteamID == 0 {
			steamID, err := GetSteamIDByIGDBIDCache(item.IGDBID)
			if err == nil {
				item.SteamID = steamID
			}
			return item, nil
		}
	}
	item, err = OrganizeGameDownloadWithSteam(0, game)
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

func OrganizeGameDownloadManually(gameID primitive.ObjectID, platform string, platformID int) (*model.GameInfo, error) {
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

func TransformSteamIDToIGDBID() {
	gameInfos, err := db.GetGameInfoWithSteamID()
	if err != nil {
		return
	}
	for _, info := range gameInfos {
		id, err := GetIGDBIDBySteamIDCache(info.SteamID)
		if err != nil {
			continue
		}
		existedInfo, err := db.GetGameInfoByPlatformID("igdb", id)
		if err == nil {
			existedInfo.GameIDs = append(existedInfo.GameIDs, info.GameIDs...)
			existedInfo.GameIDs = utils.Unique(existedInfo.GameIDs)
			_ = db.SaveGameInfo(existedInfo)
			_ = db.DeleteGameInfoByID(info.ID)
		} else {
			if err == mongo.ErrNoDocuments {
				newInfo, err := GenerateIGDBGameInfo(id)
				if err != nil {
					continue
				}
				newInfo.ID = info.ID
				newInfo.CreatedAt = info.CreatedAt
				newInfo.GameIDs = info.GameIDs
				_ = db.SaveGameInfo(newInfo)
			}
		}
	}
}

func SupplementGameInfoPlatformID() error {
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
			_ = db.SaveGameInfo(info)
		}
	}
	return nil
}
