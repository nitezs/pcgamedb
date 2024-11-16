package handler

import (
	"net/http"

	"github.com/nitezs/pcgamedb/crawler"
	"github.com/nitezs/pcgamedb/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrganizeGameItemRequest struct {
	Platform   string `form:"platform" json:"platform" binding:"required"`
	GameID     string `form:"game_id" json:"game_id" binding:"required"`
	PlatformID int    `form:"platform_id" json:"platform_id" binding:"required"`
}

type OrganizeGameItemResponse struct {
	Status   string          `json:"status"`
	Message  string          `json:"message,omitempty"`
	GameInfo *model.GameInfo `json:"game_info,omitempty"`
}

// OrganizeGameItem organizes a specific game download.
// @Summary Organize a game download
// @Description Organizes a game download based on platform and game ID
// @Tags game
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization: Bearer <api_key>"
// @Param body body OrganizeGameItemRequest true "Organize Game Download Request"
// @Success 200 {object} OrganizeGameItemResponse
// @Failure 400 {object} OrganizeGameItemResponse
// @Failure 401 {object} OrganizeGameItemResponse
// @Failure 500 {object} OrganizeGameItemResponse
// @Security BearerAuth
// @Router /game/raw/organize [post]
func OrganizeGameItemHandler(c *gin.Context) {
	var req OrganizeGameItemRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, OrganizeGameItemResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	objID, err := primitive.ObjectIDFromHex(req.GameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, OrganizeGameItemResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	info, err := crawler.OrganizeGameItemManually(objID, req.Platform, req.PlatformID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, OrganizeGameItemResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, OrganizeGameItemResponse{
		Status:   "ok",
		GameInfo: info,
	})
}
