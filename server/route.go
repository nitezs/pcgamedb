package server

import (
	"github.com/nitezs/pcgamedb/server/handler"
	"github.com/nitezs/pcgamedb/server/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	docs "github.com/nitezs/pcgamedb/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initRoute(app *gin.Engine) {
	app.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	GameInfoGroup := app.Group("/game")
	GameItemGroup := GameInfoGroup.Group("/raw")

	GameItemGroup.GET("/unorganized", handler.GetUnorganizedGameItemsHandler)
	GameItemGroup.POST("/organize", middleware.Auth(), handler.OrganizeGameItemHandler)
	GameItemGroup.GET("/id/:id", handler.GetGameItemByIDHanlder)
	GameItemGroup.GET("/name/:name", handler.GetGameItemByRawNameHandler)
	GameItemGroup.GET("/author/:author", handler.GetGameItemsByAuthorHandler)

	GameInfoGroup.GET("/search", handler.SearchGamesHandler)
	GameInfoGroup.GET("/name/:name", handler.GetGameInfosByNameHandler)
	GameInfoGroup.GET("/platform/:platform_type/:platform_id", handler.GetGameInfoByPlatformIDHandler)
	GameInfoGroup.GET("/id/:id", handler.GetGameInfoByIDHandler)
	GameInfoGroup.PUT("/update", middleware.Auth(), handler.UpdateGameInfoHandler)
	GameInfoGroup.DELETE("/id/:id", middleware.Auth(), handler.DeleteGameInfoHandler)

	app.GET("/ranking/:type", handler.GetRankingHandler)
	app.GET("/healthcheck", handler.HealthCheckHandler)
	app.GET("/author", handler.GetAllAuthorsHandler)
	app.POST("/clean", middleware.Auth(), handler.CleanGameHandler)

	docs.SwaggerInfo.BasePath = "/api"
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
