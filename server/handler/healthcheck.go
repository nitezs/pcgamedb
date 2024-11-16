package handler

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/nitezs/pcgamedb/config"
	"github.com/nitezs/pcgamedb/constant"
	"github.com/nitezs/pcgamedb/db"

	"github.com/gin-gonic/gin"
)

type HealthCheckResponse struct {
	Version            string `json:"version"`
	Status             string `json:"status"`
	Message            string `json:"message,omitempty"`
	Date               string `json:"date"`
	Uptime             string `json:"uptime"`
	Alloc              string `json:"alloc"`
	AutoCrawl          bool   `json:"auto_crawl"`
	GameItem           int64  `json:"game_download,omitempty"`
	GameInfo           int64  `json:"game_info,omitempty"`
	Unorganized        int64  `json:"unorganized,omitempty"`
	RedisAvaliable     bool   `json:"redis_avaliable"`
	OnlineFixAvaliable bool   `json:"online_fix_avaliable"`
	MegaAvaliable      bool   `json:"mega_avaliable"`
}

// HealthCheckHandler performs a health check of the service.
// @Summary Health Check
// @Description Performs a server health check and returns detailed server status including the current time, uptime, and configuration settings such as AutoCrawl.
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} HealthCheckResponse
// @Failure 500 {string} HealthCheckResponse
// @Router /healthcheck [get]
func HealthCheckHandler(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	downloadCount, _ := db.GetGameItemCount()
	infoCount, _ := db.GetGameInfoCount()
	unorganized, err := db.GetUnorganizedGameItems(-1)
	unorganizedCount := int64(0)
	if err == nil {
		unorganizedCount = int64(len(unorganized))
	}
	c.JSON(http.StatusOK, HealthCheckResponse{
		Status:             "ok",
		Version:            constant.Version,
		Date:               time.Now().Format("2006-01-02 15:04:05"),
		Uptime:             time.Since(config.Runtime.ServerStartTime).String(),
		AutoCrawl:          config.Config.Server.AutoCrawl,
		Alloc:              fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024.0/1024.0),
		GameItem:           downloadCount,
		GameInfo:           infoCount,
		Unorganized:        unorganizedCount,
		RedisAvaliable:     config.Config.RedisAvaliable,
		OnlineFixAvaliable: config.Config.OnlineFixAvaliable,
		MegaAvaliable:      config.Config.MegaAvaliable,
	})
}
