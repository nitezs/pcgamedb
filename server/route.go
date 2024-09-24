package server

import (
	"pcgamedb/server/handler"
	"pcgamedb/server/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	docs "pcgamedb/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initRoute(app *gin.Engine) {
	app.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	GameInfoGroup := app.Group("/game")
	GameDownloadGroup := GameInfoGroup.Group("/raw")

	GameDownloadGroup.GET("/unorganized", handler.GetUnorganizedGameDownloadsHandler)
	GameDownloadGroup.POST("/organize", middleware.Auth(), handler.OrganizeGameDownloadHandler)
	GameDownloadGroup.GET("/id/:id", handler.GetGameDownloadByIDHanlder)
	GameDownloadGroup.GET("/name/:name", handler.GetGameDownloadByRawNameHandler)
	GameDownloadGroup.GET("/author/:author", handler.GetGameDownloadsByAuthorHandler)

	GameInfoGroup.GET("/search", handler.SearchGamesHandler)
	GameInfoGroup.GET("/name/:name", handler.GetGameInfosByNameHandler)
	GameInfoGroup.GET("/platform/:platform_type/:platform_id", handler.GetGameInfoByPlatformIDHandler)
	GameInfoGroup.GET("/id/:id", handler.GetGameInfoByIDHandler)
	GameInfoGroup.PUT("/update", middleware.Auth(), handler.UpdateGameInfoHandler)

	app.GET("/ranking/:type", handler.GetRankingHandler)
	app.GET("/healthcheck", handler.HealthCheckHandler)
	app.GET("/author", handler.GetAllAuthorsHandler)
	app.POST("/clean", middleware.Auth(), handler.CleanGameHandler)

	docs.SwaggerInfo.BasePath = "/api"
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
