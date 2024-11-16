package handler

import (
	"net/http"

	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetGameInfoByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

type GetGameInfoByIDResponse struct {
	Status   string          `json:"status"`
	Message  string          `json:"message,omitempty"`
	GameInfo *model.GameInfo `json:"game_info,omitempty"`
}

// GetGameInfoByID retrieves game information by ID.
// @Summary Retrieve game info by ID
// @Description Retrieves details of a game by game ID
// @Tags game
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Success 200 {object} GetGameInfoByIDResponse
// @Failure 400 {object} GetGameInfoByIDResponse
// @Failure 500 {object} GetGameInfoByIDResponse
// @Router /game/id/{id} [get]
func GetGameInfoByIDHandler(c *gin.Context) {
	var req GetGameItemByIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, GetGameInfoByIDResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, GetGameInfoByIDResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	gameInfo, err := db.GetGameInfoByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, GetGameInfoByIDResponse{
				Status:  "ok",
				Message: "No results found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, GetGameInfoByIDResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	gameInfo.Games, err = db.GetGameItemsByIDs(gameInfo.GameIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetGameInfoByIDResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, GetGameInfoByIDResponse{
		Status:   "ok",
		GameInfo: gameInfo,
	})
}
