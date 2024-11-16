package crawler

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/model"
	"github.com/nitezs/pcgamedb/utils"

	"github.com/jlaffaye/ftp"
	"go.uber.org/zap"
)

const (
	ftpAddress  = "72.21.17.26:13017"
	ftpUsername = "ARMGDDNGames"
	ftpPassword = "ARMGDDNGames"
)

type GameData struct {
	NumberOfGame string `json:"Number of game"`
	AppID        string `json:"appid"`
	FolderName   string `json:"foldername"`
}

type ARMGDDNCrawler struct {
	logger zap.Logger
	conn   *ftp.ServerConn
}

// Deprecated: ARMGDDN has changed resource distribution method
func NewARMGDDNCrawler(logger *zap.Logger) *ARMGDDNCrawler {
	return &ARMGDDNCrawler{
		logger: *logger,
	}
}

func (c *ARMGDDNCrawler) connectFTP() error {
	var err error
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	c.conn, err = ftp.Dial(ftpAddress, ftp.DialWithTimeout(5*time.Second), ftp.DialWithExplicitTLS(tlsConfig))
	if err != nil {
		return err
	}
	if err = c.conn.Login(ftpUsername, ftpPassword); err != nil {
		return err
	}
	return nil
}

func (c *ARMGDDNCrawler) fetchAndParseFTPData(filePath string) ([]GameData, error) {
	r, err := c.conn.Retr(filePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var data []GameData
	if err = json.Unmarshal(buf, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *ARMGDDNCrawler) crawlGames(data []GameData, platform string, num int) ([]*model.GameItem, error) {
	count := 0
	var res []*model.GameItem
	modTimeMap := make(map[string]time.Time)
	entries, err := c.conn.List(fmt.Sprintf("/%s", platform))
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.Type == ftp.EntryTypeFolder {
			modTimeMap[entry.Name] = entry.Time
		}
	}
	for _, v := range data {
		if count == num {
			break
		}
		path := fmt.Sprintf("/%s/%s", platform, v.FolderName)
		u := fmt.Sprintf("ARMGDDNGames/%s/%s", platform, v.NumberOfGame)
		modTime, ok := modTimeMap[v.FolderName]
		if !ok {
			c.logger.Warn("mod time not found", zap.String("url", u))
			continue
		}
		updateFlag := fmt.Sprintf("ARMGDDNGames/%s/%s/%s", platform, v.NumberOfGame, modTime.UTC().String())
		if db.IsARMGDDNCrawled(updateFlag) {
			continue
		}
		c.logger.Info("Crawling", zap.String("url", u))
		walker := c.conn.Walk(path)
		size := int64(0)
		for walker.Next() {
			if walker.Stat().Type == ftp.EntryTypeFile {
				fileSize, err := c.conn.FileSize(walker.Path())
				if err != nil {
					c.logger.Warn("file size error", zap.Error(err))
					break
				}
				size += fileSize
			}
		}
		item, err := db.GetGameItemByUrl(u)
		if err != nil {
			continue
		}
		item.Url = u
		item.Name = ARMGDDNFormatter(v.FolderName)
		item.UpdateFlag = updateFlag
		item.Size = utils.FormatSize(size)
		item.RawName = v.FolderName
		item.Author = "ARMGDDN"
		item.Download = fmt.Sprintf("ftpes://%s:%s@%s/%s/%s", ftpUsername, ftpPassword, ftpAddress, platform, url.QueryEscape(v.FolderName))
		if err := db.SaveGameItem(item); err != nil {
			continue
		}
		res = append(res, item)
		count++
		var id int
		var info *model.GameInfo
		if v.AppID != "NONSTEAM" {
			id, err = strconv.Atoi(v.AppID)
			if err != nil {
				c.logger.Warn("strconv error", zap.Error(err))
				continue
			}
			info, err = OrganizeGameItemWithSteam(id, item)
			if err != nil {
				continue
			}
		} else {
			info, err = OrganizeGameItem(item)
			if err != nil {
				continue
			}
		}
		err = db.SaveGameInfo(info)
		if err != nil {
			c.logger.Warn("save game info error", zap.Error(err))
			continue
		}
	}
	return res, nil
}

func ARMGDDNFormatter(name string) string {
	cleanedName := strings.ReplaceAll(strings.TrimSpace(name), "-ARMGDDN", "")
	matchIndex := regexp.MustCompile(`v\d`).FindStringIndex(cleanedName)
	if matchIndex == nil {
		return cleanedName
	}
	return strings.TrimSpace(cleanedName[:matchIndex[0]])
}

func (c *ARMGDDNCrawler) CrawlPC(num int) ([]*model.GameItem, error) {
	return c.crawlPlatform("/PC/currentserverPC-FTP.json", "PC", num)
}

func (c *ARMGDDNCrawler) CrawlPCVR(num int) ([]*model.GameItem, error) {
	return c.crawlPlatform("/PCVR/currentserverPCVR-FTP.json", "PCVR", num)
}

func (c *ARMGDDNCrawler) Crawl(num int) ([]*model.GameItem, error) {
	num1 := num / 2
	num2 := num - num1
	if num == -1 {
		num1 = -1
		num2 = -1
	}
	res1, err := c.CrawlPC(num1)
	if err != nil {
		return nil, err
	}
	res2, err := c.CrawlPCVR(num2)
	if err != nil {
		return nil, err
	}
	return append(res1, res2...), nil
}

func (c *ARMGDDNCrawler) CrawlAll() ([]*model.GameItem, error) {
	return c.Crawl(-1)
}

func (c *ARMGDDNCrawler) crawlPlatform(jsonFile, platform string, num int) ([]*model.GameItem, error) {
	err := c.connectFTP()
	if err != nil {
		return nil, err
	}
	defer func() { _ = c.conn.Quit() }()

	data, err := c.fetchAndParseFTPData(jsonFile)
	if err != nil {
		return nil, err
	}

	return c.crawlGames(data, platform, num)
}
