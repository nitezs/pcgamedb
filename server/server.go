package server

import (
	"io"
	"time"

	"github.com/nitezs/pcgamedb/cache"
	"github.com/nitezs/pcgamedb/config"
	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/log"
	"github.com/nitezs/pcgamedb/server/middleware"
	"github.com/nitezs/pcgamedb/task"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func init() {
	config.Runtime.ServerStartTime = time.Now()
}

func Run() {
	db.CheckConnect()
	cache.CheckConnect()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	app := gin.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recovery())
	initRoute(app)
	log.Logger.Info("Server running", zap.String("port", config.Config.Server.Port))
	if config.Config.Server.AutoCrawl {
		go func() {
			c := cron.New()
			_, err := c.AddFunc("0 */3 * * *", func() { task.Crawl(log.TaskLogger) })
			if err != nil {
				log.Logger.Error("Error adding cron job", zap.Error(err))
			}
			c.Start()
		}()
	}
	err := app.Run(":" + config.Config.Server.Port)
	if err != nil {
		log.Logger.Panic("Failed to run server", zap.Error(err))
	}
}
