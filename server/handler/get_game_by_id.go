package handler

import (
	"net/http"

	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetGameItemByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

type GetGameItemByIDResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message,omitempty"`
	Game    *model.GameItem `json:"game,omitempty"`
}

// GetGameItemByID retrieves game download details by ID.
// @Summary Retrieve game download by ID
// @Description Retrieves details of a game download by game ID
// @Tags game
// @Accept json
// @Produce json
// @Param id path string true "Game Download ID"
// @Success 200 {object} GetGameItemByIDResponse
// @Failure 400 {object} GetGameItemByIDResponse
// @Failure 500 {object} GetGameItemByIDResponse
// @Router /game/raw/id/{id} [get]
func GetGameItemByIDHanlder(c *gin.Context) {
	var req GetGameItemByIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, GetGameItemByIDResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, GetGameItemByIDResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	game, err := db.GetGameItemByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, GetGameItemByIDResponse{
				Status:  "ok",
				Message: "No results found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, GetGameItemByIDResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, GetGameItemByIDResponse{
		Status: "ok",
		Game:   game,
	})
}
