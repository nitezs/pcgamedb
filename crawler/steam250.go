package crawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/nitezs/pcgamedb/cache"
	"github.com/nitezs/pcgamedb/config"
	"github.com/nitezs/pcgamedb/constant"
	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/model"
	"github.com/nitezs/pcgamedb/utils"

	"github.com/PuerkitoBio/goquery"
)

func GetSteam250(url string) ([]*model.GameInfo, error) {
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: url,
	})
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		return nil, err
	}
	var rank []model.Steam250Item
	var item model.Steam250Item
	steamIDs := make([]int, 0)
	doc.Find(".appline").Each(func(i int, s *goquery.Selection) {
		item.Name = s.Find(".title>a").First().Text()
		idStr := s.Find(".store").AttrOr("href", "")
		idSlice := regexp.MustCompile(`app/(\d+)/`).FindStringSubmatch(idStr)
		if len(idSlice) < 2 {
			return
		}
		item.SteamID, _ = strconv.Atoi(idSlice[1])
		rank = append(rank, item)
		steamIDs = append(steamIDs, item.SteamID)
	})
	var res []*model.GameInfo
	count := 0
	idMap, err := GetIGDBIDBySteamIDsCache(steamIDs)
	if err != nil {
		return nil, err
	}
	for _, item := range rank {
		if count == 10 {
			break
		}
		if idMap[item.SteamID] != 0 {
			info, err := db.GetGameInfoByPlatformID("igdb", idMap[item.SteamID])
			if err == nil {
				res = append(res, info)
				count++
				continue
			}
		} else {
			info, err := db.GetGameInfoByPlatformID("steam", item.SteamID)
			if err == nil {
				res = append(res, info)
				count++
				continue
			}
		}
	}
	return res, nil
}

func GetSteam250Top250() ([]*model.GameInfo, error) {
	return GetSteam250(constant.Steam250Top250URL)
}

func GetSteam250Top250Cache() ([]*model.GameInfo, error) {
	return GetSteam250Cache("top250", GetSteam250Top250)
}

func GetSteam250BestOfTheYear() ([]*model.GameInfo, error) {
	return GetSteam250(fmt.Sprintf(constant.Steam250BestOfTheYearURL, time.Now().UTC().Year()))
}

func GetSteam250BestOfTheYearCache() ([]*model.GameInfo, error) {
	return GetSteam250Cache(fmt.Sprintf("bestoftheyear:%v", time.Now().UTC().Year()), GetSteam250BestOfTheYear)
}

func GetSteam250WeekTop50() ([]*model.GameInfo, error) {
	return GetSteam250(constant.Steam250WeekTop50URL)
}

func GetSteam250WeekTop50Cache() ([]*model.GameInfo, error) {
	return GetSteam250Cache("weektop50", GetSteam250WeekTop50)
}

func GetSteam250MostPlayed() ([]*model.GameInfo, error) {
	return GetSteam250(constant.Steam250MostPlayedURL)
}

func GetSteam250MostPlayedCache() ([]*model.GameInfo, error) {
	return GetSteam250Cache("mostplayed", GetSteam250MostPlayed)
}

func GetSteam250Cache(k string, f func() ([]*model.GameInfo, error)) ([]*model.GameInfo, error) {
	if config.Config.RedisAvaliable {
		key := k
		val, exist := cache.Get(key)
		if exist {
			var res []*model.GameInfo
			err := json.Unmarshal([]byte(val), &res)
			if err != nil {
				return nil, err
			}
			return res, nil
		} else {
			data, err := f()
			if err != nil {
				return nil, err
			}
			dataBytes, err := json.Marshal(data)
			if err != nil {
				return data, nil
			}
			err = cache.AddWithExpire(key, dataBytes, 24*time.Hour)
			if err != nil {
				return data, nil
			}
			return data, nil
		}
	} else {
		return f()
	}
}
