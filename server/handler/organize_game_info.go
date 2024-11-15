package handler

import (
	"net/http"

	"github.com/nitezs/pcgamedb/crawler"
	"github.com/nitezs/pcgamedb/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrganizeGameDownloadRequest struct {
	Platform   string `form:"platform" json:"platform" binding:"required"`
	GameID     string `form:"game_id" json:"game_id" binding:"required"`
	PlatformID int    `form:"platform_id" json:"platform_id" binding:"required"`
}

type OrganizeGameDownloadResponse struct {
	Status   string          `json:"status"`
	Message  string          `json:"message,omitempty"`
	GameInfo *model.GameInfo `json:"game_info,omitempty"`
}

// OrganizeGameDownload organizes a specific game download.
// @Summary Organize a game download
// @Description Organizes a game download based on platform and game ID
// @Tags game
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization: Bearer <api_key>"
// @Param body body OrganizeGameDownloadRequest true "Organize Game Download Request"
// @Success 200 {object} OrganizeGameDownloadResponse
// @Failure 400 {object} OrganizeGameDownloadResponse
// @Failure 401 {object} OrganizeGameDownloadResponse
// @Failure 500 {object} OrganizeGameDownloadResponse
// @Security BearerAuth
// @Router /game/raw/organize [post]
func OrganizeGameDownloadHandler(c *gin.Context) {
	var req OrganizeGameDownloadRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, OrganizeGameDownloadResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	objID, err := primitive.ObjectIDFromHex(req.GameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, OrganizeGameDownloadResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	info, err := crawler.OrganizeGameDownloadManually(objID, req.Platform, req.PlatformID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, OrganizeGameDownloadResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, OrganizeGameDownloadResponse{
		Status:   "ok",
		GameInfo: info,
	})
}
