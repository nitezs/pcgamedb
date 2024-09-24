package handler

import (
	"net/http"
	"pcgamedb/db"
	"pcgamedb/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetGameDownloadByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

type GetGameDownloadByIDResponse struct {
	Status  string              `json:"status"`
	Message string              `json:"message,omitempty"`
	Game    *model.GameDownload `json:"game,omitempty"`
}

// GetGameDownloadByID retrieves game download details by ID.
// @Summary Retrieve game download by ID
// @Description Retrieves details of a game download by game ID
// @Tags game
// @Accept json
// @Produce json
// @Param id path string true "Game Download ID"
// @Success 200 {object} GetGameDownloadByIDResponse
// @Failure 400 {object} GetGameDownloadByIDResponse
// @Failure 500 {object} GetGameDownloadByIDResponse
// @Router /game/raw/id/{id} [get]
func GetGameDownloadByIDHanlder(c *gin.Context) {
	var req GetGameDownloadByIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, GetGameDownloadByIDResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, GetGameDownloadByIDResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	game, err := db.GetGameDownloadByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, GetGameDownloadByIDResponse{
				Status:  "ok",
				Message: "No results found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, GetGameDownloadByIDResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, GetGameDownloadByIDResponse{
		Status: "ok",
		Game:   game,
	})
}
