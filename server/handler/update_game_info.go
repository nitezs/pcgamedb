package handler

import (
	"net/http"
	"pcgamedb/crawler"
	"pcgamedb/db"
	"pcgamedb/model"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateGameInfoRequest struct {
	GameID     string `json:"game_id" binding:"required"`
	Platform   string `json:"platform" binding:"required"`
	PlatformID int    `json:"platform_id" binding:"required"`
}

type UpdateGameInfoResponse struct {
	Status   string          `json:"status"`
	Message  string          `json:"message"`
	GameInfo *model.GameInfo `json:"game_info,omitempty"`
}

// UpdateGameInfoHandler updates game information.
// @Summary Update game info
// @Description Updates details of a game
// @Tags game
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization: Bearer <api_key>"
// @Param body body handler.UpdateGameInfoRequest true "Update Game Info Request"
// @Success 200 {object} handler.UpdateGameInfoResponse
// @Failure 400 {object} handler.UpdateGameInfoResponse
// @Failure 401 {object} handler.UpdateGameInfoResponse
// @Failure 500 {object} handler.UpdateGameInfoResponse
// @Router /game/update [post]
func UpdateGameInfoHandler(c *gin.Context) {
	var req UpdateGameInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, UpdateGameInfoResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	req.Platform = strings.ToLower(req.Platform)
	platformMap := map[string]bool{
		"steam": true,
		"igdb":  true,
		"gog":   true,
	}
	if _, ok := platformMap[req.Platform]; !ok {
		c.JSON(http.StatusBadRequest, UpdateGameInfoResponse{
			Status:  "error",
			Message: "Invalid platform",
		})
		return
	}
	objID, err := primitive.ObjectIDFromHex(req.GameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, UpdateGameInfoResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	info, err := db.GetGameInfoByID(objID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateGameInfoResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	newInfo, err := crawler.GenerateGameInfo(req.Platform, req.PlatformID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateGameInfoResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	newInfo.ID = objID
	newInfo.GameIDs = info.GameIDs
	err = db.SaveGameInfo(newInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateGameInfoResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UpdateGameInfoResponse{
		Status:   "ok",
		Message:  "Game info updated successfully",
		GameInfo: newInfo,
	})
}
