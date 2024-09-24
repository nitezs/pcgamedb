package server

import (
	"io"
	"pcgamedb/cache"
	"pcgamedb/config"
	"pcgamedb/db"
	"pcgamedb/log"
	"pcgamedb/server/middleware"
	"pcgamedb/task"
	"time"

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
