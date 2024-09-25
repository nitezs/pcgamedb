package handler

import (
	"net/http"
	"pcgamedb/db"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteGameInfoRequest struct {
	ID string `uri:"id" binding:"required"`
}

type DeleteGameInfoResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// DeleteGameInfoHandler is the handler for deleting game info
// @Summary Delete game info by ID
// @Description Delete game info by ID
// @Tags game
// @Produce json
// @Param Authorization header string true "Authorization: Bearer <api_key>"
// @Param id path string true "Game ID"
// @Success 200 {object} DeleteGameInfoResponse
// @Failure 400 {object} DeleteGameInfoResponse
// @Failure 500 {object} DeleteGameInfoResponse
// @Security BearerAuth
// @Router /game/id/{id} [delete]
func DeleteGameInfoHandler(c *gin.Context) {
	var req DeleteGameInfoRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, DeleteGameInfoResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	objID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, DeleteGameInfoResponse{
			Status:  "error",
			Message: "Invalid ID",
		})
		return
	}
	err = db.DeleteGameInfoByID(objID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, DeleteGameInfoResponse{
			Status:  "error",
			Message: "Failed to delete game",
		})
		return
	}
	c.JSON(http.StatusOK, DeleteGameInfoResponse{
		Status:  "success",
		Message: "Game deleted successfully",
	})
}
