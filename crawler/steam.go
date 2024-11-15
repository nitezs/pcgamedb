package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/nitezs/pcgamedb/cache"
	"github.com/nitezs/pcgamedb/config"
	"github.com/nitezs/pcgamedb/constant"
	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/model"
	"github.com/nitezs/pcgamedb/utils"
)

func _GetSteamID(name string) (int, error) {
	baseURL, _ := url.Parse(constant.SteamSearchURL)
	params := url.Values{}
	params.Add("term", name)
	baseURL.RawQuery = params.Encode()

	resp, err := utils.Fetch(utils.FetchConfig{
		Url: baseURL.String(),
	})
	if err != nil {
		return 0, err
	}
	idRegex := regexp.MustCompile(`data-ds-appid="(.*?)"`)
	nameRegex := regexp.MustCompile(`<span class="title">(.*?)</span>`)
	idRegexRes := idRegex.FindAllStringSubmatch(string(resp.Data), -1)
	nameRegexRes := nameRegex.FindAllStringSubmatch(string(resp.Data), -1)

	if len(idRegexRes) == 0 {
		return 0, fmt.Errorf("Steam ID not found: %s", name)
	}

	maxSim := 0.0
	maxSimID := 0
	for i, id := range idRegexRes {
		idStr := id[1]
		nameStr := nameRegexRes[i][1]
		if index := strings.Index(idStr, ","); index != -1 {
			idStr = idStr[:index]
		}
		if strings.EqualFold(strings.TrimSpace(nameStr), strings.TrimSpace(name)) {
			return strconv.Atoi(idStr)
		} else {
			sim := utils.Similarity(nameStr, name)
			if sim >= 0.8 && sim > maxSim {
				maxSim = sim
				maxSimID, _ = strconv.Atoi(idStr)
			}
		}
	}
	if maxSimID != 0 {
		return maxSimID, nil
	}
	return 0, fmt.Errorf("Steam ID not found: %s", name)
}

func GetSteamID(name string) (int, error) {
	name1 := name
	name2 := FormatName(name)
	names := []string{name1}
	if name1 != name2 {
		names = append(names, name2)
	}
	for _, n := range names {
		id, err := _GetSteamID(n)
		if err == nil {
			return id, nil
		}
	}
	return 0, errors.New("Steam ID not found")
}

func GetSteamIDCache(name string) (int, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("steam_id:%s", name)
		val, exist := cache.Get(key)
		if exist {
			id, err := strconv.Atoi(val)
			if err != nil {
				return 0, err
			}
			return id, nil
		} else {
			id, err := GetSteamID(name)
			if err != nil {
				return 0, err
			}
			_ = cache.Add(key, id)
			return id, nil
		}
	} else {
		return GetSteamID(name)
	}
}

func GetSteamAppDetail(id int) (*model.SteamAppDetail, error) {
	baseURL, _ := url.Parse(constant.SteamAppDetailURL)
	params := url.Values{}
	params.Add("appids", strconv.Itoa(id))
	// params.Add("l", "schinese")
	baseURL.RawQuery = params.Encode()
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: baseURL.String(),
		Headers: map[string]string{
			"User-Agent": "",
		},
	})
	if err != nil {
		return nil, err
	}
	var detail map[string]*model.SteamAppDetail
	if err = json.Unmarshal(resp.Data, &detail); err != nil {
		return nil, err
	}
	if _, ok := detail[strconv.Itoa(id)]; !ok {
		return nil, fmt.Errorf("Steam App not found: %d", id)
	}
	if detail[strconv.Itoa(id)] == nil {
		return nil, fmt.Errorf("Steam App not found: %d", id)
	}
	return detail[strconv.Itoa(id)], nil
}

func GetSteamAppDetailCache(id int) (*model.SteamAppDetail, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("steam_game:%d", id)
		val, exist := cache.Get(key)
		if exist {
			var detail model.SteamAppDetail
			if err := json.Unmarshal([]byte(val), &detail); err != nil {
				return nil, err
			}
			return &detail, nil
		} else {
			data, err := GetSteamAppDetail(id)
			if err != nil {
				return nil, err
			}
			dataBytes, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}
			_ = cache.Add(key, dataBytes)
			return data, nil
		}
	} else {
		return GetSteamAppDetail(id)
	}
}

func GenerateSteamGameInfo(id int) (*model.GameInfo, error) {
	item := &model.GameInfo{}
	detail, err := GetSteamAppDetailCache(id)
	if err != nil {
		return nil, err
	}
	item.SteamID = id
	item.Name = detail.Data.Name
	item.Description = detail.Data.ShortDescription
	item.Cover = fmt.Sprintf("https://shared.cloudflare.steamstatic.com/store_item_assets/steam/apps/%v/library_600x900_2x.jpg", id)
	item.Developers = detail.Data.Developers
	item.Publishers = detail.Data.Publishers
	screenshots := []string{}
	for _, screenshot := range detail.Data.Screenshots {
		screenshots = append(screenshots, screenshot.PathFull)
	}
	item.Screenshots = screenshots
	return item, nil
}

func OrganizeGameDownloadWithSteam(id int, game *model.GameDownload) (*model.GameInfo, error) {
	var err error
	if id == 0 {
		id, err = GetSteamIDCache(game.Name)
		if err != nil {
			return nil, err
		}
	}
	d, err := db.GetGameInfoByPlatformID("steam", id)
	if err == nil {
		d.GameIDs = append(d.GameIDs, game.ID)
		d.GameIDs = utils.Unique(d.GameIDs)
		return d, nil
	}
	detail, err := GenerateGameInfo("steam", id)
	if err != nil {
		return nil, err
	}
	detail.GameIDs = append(detail.GameIDs, game.ID)
	detail.GameIDs = utils.Unique(detail.GameIDs)
	return detail, nil
}

func GetSteamIDByIGDBID(IGDBID int) (int, error) {
	var err error
	if TwitchToken == "" {
		TwitchToken, err = LoginTwitch()
		if err != nil {
			return 0, err
		}
	}
	resp, err := utils.Fetch(utils.FetchConfig{
		Url:    constant.IGDBWebsitesURL,
		Method: "POST",
		Headers: map[string]string{
			"Client-ID":     config.Config.Twitch.ClientID,
			"Authorization": "Bearer " + TwitchToken,
			"User-Agent":    "",
			"Content-Type":  "text/plain",
		},
		Data: fmt.Sprintf(`where game = %v; fields *; limit 500;`, IGDBID),
	})
	if err != nil {
		return 0, err
	}
	var data []struct {
		Game int    `json:"game"`
		Url  string `json:"url"`
	}
	if err = json.Unmarshal(resp.Data, &data); err != nil {
		return 0, err
	}
	if len(data) == 0 {
		return 0, errors.New("Not found")
	}
	for _, v := range data {
		if strings.HasPrefix(v.Url, "https://store.steampowered.com/app/") {
			regex := regexp.MustCompile(`https://store.steampowered.com/app/(\d+)/?`)
			idStr := regex.FindStringSubmatch(v.Url)
			if len(idStr) < 2 {
				return 0, errors.New("Failed parse")
			}
			steamID, err := strconv.Atoi(idStr[1])
			if err != nil {
				return 0, err
			}
			return steamID, nil
		}
	}
	return 0, errors.New("Not found")
}

func GetSteamIDByIGDBIDCache(IGDBID int) (int, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("steam_game:%d", IGDBID)
		val, exist := cache.Get(key)
		if exist {
			id, err := strconv.Atoi(val)
			if err != nil {
				return 0, err
			}
			return id, nil
		} else {
			id, err := GetSteamIDByIGDBID(IGDBID)
			if err != nil {
				return 0, err
			}
			dataBytes := strconv.Itoa(id)
			_ = cache.Add(key, dataBytes)
			return id, nil
		}
	} else {
		return GetSteamIDByIGDBID(IGDBID)
	}
}
