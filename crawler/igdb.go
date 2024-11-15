package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/nitezs/pcgamedb/cache"
	"github.com/nitezs/pcgamedb/config"
	"github.com/nitezs/pcgamedb/constant"
	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/model"
	"github.com/nitezs/pcgamedb/utils"
)

var TwitchToken string

func _GetIGDBID(name string) (int, error) {
	var err error
	if TwitchToken == "" {
		TwitchToken, err = LoginTwitch()
		if err != nil {
			return 0, fmt.Errorf("failed to login twitch: %w", err)
		}
	}
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.IGDBSearchURL,
		Headers: map[string]string{
			"Client-ID":     config.Config.Twitch.ClientID,
			"Authorization": "Bearer " + TwitchToken,
			"User-Agent":    "",
			"Content-Type":  "text/plain",
		},
		Data:   fmt.Sprintf(`search "%s"; fields *; limit 50; where game.platforms = [6] | game.platforms=[130] | game.platforms=[384] | game.platforms=[163];`, name),
		Method: "POST",
	})
	if string(resp.Data) == "[]" {
		resp, err = utils.Fetch(utils.FetchConfig{
			Url: constant.IGDBSearchURL,
			Headers: map[string]string{
				"Client-ID":     config.Config.Twitch.ClientID,
				"Authorization": "Bearer " + TwitchToken,
				"User-Agent":    "",
				"Content-Type":  "text/plain",
			},
			Data:   fmt.Sprintf(`search "%s"; fields *; limit 50;`, name),
			Method: "POST",
		})
	}
	if err != nil {
		return 0, err
	}
	var data model.IGDBSearches
	if err = json.Unmarshal(resp.Data, &data); err != nil {
		return 0, fmt.Errorf("failed to unmarshal: %w, %s", err, debug.Stack())
	}
	if len(data) == 1 {
		return data[0].Game, nil
	}
	for _, item := range data {
		if strings.EqualFold(item.Name, name) {
			return item.Game, nil
		}
		if utils.Similarity(name, item.Name) >= 0.8 {
			return item.Game, nil
		}
		detail, err := GetIGDBAppDetailCache(item.Game)
		if err != nil {
			return 0, err
		}
		for _, alternativeNames := range detail.AlternativeNames {
			if utils.Similarity(alternativeNames.Name, name) >= 0.8 {
				return item.Game, nil
			}
		}
	}
	return 0, fmt.Errorf("IGDB ID not found: %s", name)
}

func GetIGDBID(name string) (int, error) {
	name1 := name
	name2 := FormatName(name)
	names := []string{name1}
	if name1 != name2 {
		names = append(names, name2)
	}
	for _, name := range names {
		id, err := _GetIGDBID(name)
		if err == nil {
			return id, nil
		}
	}
	return 0, errors.New("IGDB ID not found")
}

func GetIGDBIDCache(name string) (int, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("igdb_id:%s", name)
		val, exist := cache.Get(key)
		if exist {
			id, err := strconv.Atoi(val)
			if err != nil {
				return 0, err
			}
			return id, nil
		} else {
			id, err := GetIGDBID(name)
			if err != nil {
				return 0, err
			}
			_ = cache.Add(key, id)
			return id, nil
		}
	} else {
		return GetIGDBID(name)
	}
}

func GetIGDBAppDetail(id int) (*model.IGDBGameDetail, error) {
	var err error
	if TwitchToken == "" {
		TwitchToken, err = LoginTwitch()
		if err != nil {
			return nil, err
		}
	}
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.IGDBGameURL,
		Headers: map[string]string{
			"Client-ID":     config.Config.Twitch.ClientID,
			"Authorization": "Bearer " + TwitchToken,
			"User-Agent":    "",
			"Content-Type":  "text/plain",
		},
		Data:   fmt.Sprintf(`where id=%v ;fields *,alternative_names.name,language_supports.language,language_supports.language_support_type,screenshots.url,cover.url,involved_companies.company,involved_companies.developer,involved_companies.publisher;`, id),
		Method: "POST",
	})
	if err != nil {
		return nil, err
	}
	var data model.IGDBGameDetails
	if err = json.Unmarshal(resp.Data, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("IGDB App not found")
	}
	if data[0].Name == "" {
		return GetIGDBAppDetail(id)
	}
	return data[0], nil
}

func GetIGDBAppDetailCache(id int) (*model.IGDBGameDetail, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("igdb_game:%v", id)
		val, exist := cache.Get(key)
		if exist {
			var data model.IGDBGameDetail
			if err := json.Unmarshal([]byte(val), &data); err != nil {
				return nil, err
			}
			return &data, nil
		} else {
			data, err := GetIGDBAppDetail(id)
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
		return GetIGDBAppDetail(id)
	}
}

func LoginTwitch() (string, error) {
	baseURL, _ := url.Parse(constant.TwitchAuthURL)
	params := url.Values{}
	params.Add("client_id", config.Config.Twitch.ClientID)
	params.Add("client_secret", config.Config.Twitch.ClientSecret)
	params.Add("grant_type", "client_credentials")
	baseURL.RawQuery = params.Encode()
	resp, err := utils.Fetch(utils.FetchConfig{
		Url:    baseURL.String(),
		Method: "POST",
		Headers: map[string]string{
			"User-Agent": "",
		},
	})
	if err != nil {
		return "", err
	}
	data := struct {
		AccessToken string `json:"access_token"`
	}{}
	err = json.Unmarshal(resp.Data, &data)
	if err != nil {
		return "", err
	}
	return data.AccessToken, nil
}

func GetIGDBCompany(id int) (string, error) {
	var err error
	if TwitchToken == "" {
		TwitchToken, err = LoginTwitch()
		if err != nil {
			return "", err
		}
	}
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.IGDBCompaniesURL,
		Headers: map[string]string{
			"Client-ID":     config.Config.Twitch.ClientID,
			"Authorization": "Bearer " + TwitchToken,
			"User-Agent":    "",
			"Content-Type":  "text/plain",
		},
		Data:   fmt.Sprintf(`where id=%v; fields *;`, id),
		Method: "POST",
	})
	if err != nil {
		return "", err
	}
	var data model.IGDBCompanies
	if err = json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", errors.New("Not found")
	}
	if data[0].Name == "" {
		return GetIGDBCompany(id)
	}
	return data[0].Name, nil
}

func GetIGDBCompanyCache(id int) (string, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("igdb_companies:%v", id)
		val, exist := cache.Get(key)
		if exist {
			return val, nil
		} else {
			data, err := GetIGDBCompany(id)
			if err != nil {
				return "", err
			}
			_ = cache.Add(key, data)
			return data, nil
		}
	} else {
		return GetIGDBCompany(id)
	}
}

func GenerateIGDBGameInfo(id int) (*model.GameInfo, error) {
	item := &model.GameInfo{}
	detail, err := GetIGDBAppDetailCache(id)
	if err != nil {
		return nil, err
	}
	item.IGDBID = id
	item.Name = detail.Name
	item.Description = detail.Summary
	item.Cover = strings.Replace(detail.Cover.URL, "t_thumb", "t_original", 1)

	for _, lang := range detail.LanguageSupports {
		if lang.LanguageSupportType == 3 {
			l, exist := constant.IGDBLanguages[lang.Language]
			if !exist {
				continue
			}
			item.Languages = append(item.Languages, l.Name)
		}
	}

	for _, screenshot := range detail.Screenshots {
		item.Screenshots = append(item.Screenshots, strings.Replace(screenshot.URL, "t_thumb", "t_original", 1))
	}

	for _, alias := range detail.AlternativeNames {
		item.Aliases = append(item.Aliases, alias.Name)
	}

	for _, company := range detail.InvolvedCompanies {
		if company.Developer || company.Publisher {
			companyName, err := GetIGDBCompanyCache(company.Company)
			if err != nil {
				continue
			}
			if company.Developer {
				item.Developers = append(item.Developers, companyName)
			}
			if company.Publisher {
				item.Publishers = append(item.Publishers, companyName)
			}
		}
	}

	return item, nil
}

func OrganizeGameDownloadWithIGDB(id int, game *model.GameDownload) (*model.GameInfo, error) {
	var err error
	if id == 0 {
		id, err = GetIGDBIDCache(game.Name)
		if err != nil {
			return nil, err
		}
	}
	d, err := db.GetGameInfoByPlatformID("igdb", id)
	if err == nil {
		d.GameIDs = append(d.GameIDs, game.ID)
		d.GameIDs = utils.Unique(d.GameIDs)
		return d, nil
	}
	info, err := GenerateGameInfo("igdb", id)
	if err != nil {
		return nil, err
	}
	info.GameIDs = append(info.GameIDs, game.ID)
	info.GameIDs = utils.Unique(info.GameIDs)
	return info, nil
}

func GetIGDBIDBySteamID(id int) (int, error) {
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
		Data: fmt.Sprintf(`where url = "https://store.steampowered.com/app/%v" | url = "https://store.steampowered.com/app/%v/"*; fields *; limit 500;`, id, id),
	})
	if err != nil {
		return 0, err
	}
	var data []struct {
		Game int `json:"game"`
	}
	if err = json.Unmarshal(resp.Data, &data); err != nil {
		return 0, err
	}
	if len(data) == 0 {
		return 0, errors.New("Not found")
	}
	if data[0].Game == 0 {
		return GetIGDBIDBySteamID(id)
	}
	return data[0].Game, nil
}

func GetIGDBIDBySteamIDCache(id int) (int, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("igdb_id_by_steam_id:%v", id)
		val, exist := cache.Get(key)
		if exist {
			return strconv.Atoi(val)
		} else {
			data, err := GetIGDBIDBySteamID(id)
			if err != nil {
				return 0, err
			}
			_ = cache.Add(key, strconv.Itoa(data))
			return data, nil
		}
	} else {
		return GetIGDBIDBySteamID(id)
	}
}

func GetIGDBIDBySteamIDs(ids []int) (map[int]int, error) {
	var err error
	if TwitchToken == "" {
		TwitchToken, err = LoginTwitch()
		if err != nil {
			return nil, err
		}
	}
	conditionBuilder := strings.Builder{}
	for _, id := range ids {
		conditionBuilder.WriteString(fmt.Sprintf(`url = "https://store.steampowered.com/app/%v" | `, id))
		conditionBuilder.WriteString(fmt.Sprintf(`url = "https://store.steampowered.com/app/%v/"* | `, id))
	}
	condition := strings.TrimSuffix(conditionBuilder.String(), " | ")
	respBody := fmt.Sprintf(`where %s; fields *; limit 500;`, condition)
	resp, err := utils.Fetch(utils.FetchConfig{
		Url:    constant.IGDBWebsitesURL,
		Method: "POST",
		Headers: map[string]string{
			"Client-ID":     config.Config.Twitch.ClientID,
			"Authorization": "Bearer " + TwitchToken,
			"User-Agent":    "",
			"Content-Type":  "text/plain",
		},
		Data: respBody,
	})
	if err != nil {
		return nil, err
	}
	var data []struct {
		Game int    `json:"game"`
		Url  string `json:"url"`
	}
	if err = json.Unmarshal(resp.Data, &data); err != nil {
		return nil, err
	}
	ret := make(map[int]int)
	regex := regexp.MustCompile(`https://store.steampowered.com/app/(\d+)/?`)
	for _, d := range data {
		idStr := regex.FindStringSubmatch(d.Url)
		if len(idStr) < 2 {
			continue
		}
		id, err := strconv.Atoi(idStr[1])
		if err == nil {
			ret[id] = d.Game
		}
	}
	for _, id := range ids {
		if _, ok := ret[id]; !ok {
			ret[id] = 0
		}
	}
	return ret, nil
}

func GetIGDBIDBySteamIDsCache(ids []int) (map[int]int, error) {
	res := make(map[int]int)
	notExistIDs := make([]int, 0)
	if config.Config.RedisAvaliable {
		for _, steamID := range ids {
			key := fmt.Sprintf("igdb_id_by_steam_id:%v", steamID)
			val, exist := cache.Get(key)
			if exist {
				igdbID, _ := strconv.Atoi(val)
				res[steamID] = igdbID
			} else {
				notExistIDs = append(notExistIDs, steamID)
			}
		}
		if len(res) == len(ids) {
			return res, nil
		}
		idMap, err := GetIGDBIDBySteamIDs(notExistIDs)
		if err != nil {
			return nil, err
		}
		for steamID, igdbID := range idMap {
			res[steamID] = igdbID
			if igdbID != 0 {
				_ = cache.Add(fmt.Sprintf("igdb_id_by_steam_id:%v", steamID), igdbID)
			}
		}
		return res, nil
	} else {
		return GetIGDBIDBySteamIDs(ids)
	}
}
