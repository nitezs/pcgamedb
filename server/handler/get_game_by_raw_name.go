package handler

import (
	"net/http"

	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetGameItemByRawNameRequest struct {
	Name string `uri:"name" binding:"required"`
}

type GetGameItemByRawNameResponse struct {
	Status   string            `json:"status"`
	Message  string            `json:"message,omitempty"`
	GameItem []*model.GameItem `json:"game_downloads,omitempty"`
}

// GetGameItemByRawName retrieves game download details by raw name.
// @Summary Retrieve game download by raw name
// @Description Retrieves details of a game download by its raw name
// @Tags game
// @Accept json
// @Produce json
// @Param name path string true "Game Download Raw Name"
// @Success 200 {object} GetGameItemByRawNameResponse
// @Failure 400 {object} GetGameItemByRawNameResponse
// @Failure 500 {object} GetGameItemByRawNameResponse
// @Router /game/raw/name/{name} [get]
func GetGameItemByRawNameHandler(c *gin.Context) {
	var req GetGameItemByRawNameRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, GetGameItemByRawNameResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	gameDownload, err := db.GetGameItemByRawName(req.Name)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, GetGameItemByRawNameResponse{
				Status:  "ok",
				Message: "No results found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, GetGameItemByRawNameResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	if gameDownload == nil {
		c.JSON(http.StatusOK, GetGameItemByRawNameResponse{
			Status:  "ok",
			Message: "No results found",
		})
		return
	}
	c.JSON(http.StatusOK, GetGameItemByRawNameResponse{
		Status:   "ok",
		GameItem: gameDownload,
	})
}
