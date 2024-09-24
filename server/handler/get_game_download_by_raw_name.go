package handler

import (
	"net/http"
	"pcgamedb/db"
	"pcgamedb/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetGameDownloadByRawNameRequest struct {
	Name string `uri:"name" binding:"required"`
}

type GetGameDownloadByRawNameResponse struct {
	Status       string                `json:"status"`
	Message      string                `json:"message,omitempty"`
	GameDownload []*model.GameDownload `json:"game_downloads,omitempty"`
}

// GetGameDownloadByRawName retrieves game download details by raw name.
// @Summary Retrieve game download by raw name
// @Description Retrieves details of a game download by its raw name
// @Tags game
// @Accept json
// @Produce json
// @Param name path string true "Game Download Raw Name"
// @Success 200 {object} GetGameDownloadByRawNameResponse
// @Failure 400 {object} GetGameDownloadByRawNameResponse
// @Failure 500 {object} GetGameDownloadByRawNameResponse
// @Router /game/raw/name/{name} [get]
func GetGameDownloadByRawNameHandler(c *gin.Context) {
	var req GetGameDownloadByRawNameRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, GetGameDownloadByRawNameResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	gameDownload, err := db.GetGameDownloadByRawName(req.Name)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, GetGameDownloadByRawNameResponse{
				Status:  "ok",
				Message: "No results found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, GetGameDownloadByRawNameResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	if gameDownload == nil {
		c.JSON(http.StatusOK, GetGameDownloadByRawNameResponse{
			Status:  "ok",
			Message: "No results found",
		})
		return
	}
	c.JSON(http.StatusOK, GetGameDownloadByRawNameResponse{
		Status:       "ok",
		GameDownload: gameDownload,
	})
}
