package handler

import (
	"net/http"

	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetGameInfoByPlatformIDRequest struct {
	PlatformType string `uri:"platform_type" binding:"required"`
	PlatformID   int    `uri:"platform_id" binding:"required"`
}

type GetGameInfoByPlatformIDResponse struct {
	Status   string          `json:"status"`
	Message  string          `json:"message,omitempty"`
	GameInfo *model.GameInfo `json:"game_info,omitempty"`
}

// GetGameInfoByPlatformID retrieves game information by platform and ID.
// @Summary Retrieve game info by platform ID
// @Description Retrieves game information based on a platform type and platform ID
// @Tags game
// @Accept json
// @Produce json
// @Param platform_type path string true "Platform Type"
// @Param platform_id path int true "Platform ID"
// @Success 200 {object} GetGameInfoByPlatformIDResponse
// @Failure 400 {object} GetGameInfoByPlatformIDResponse
// @Failure 500 {object} GetGameInfoByPlatformIDResponse
// @Router /game/platform/{platform_type}/{platform_id} [get]
func GetGameInfoByPlatformIDHandler(c *gin.Context) {
	var req GetGameInfoByPlatformIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, GetGameInfoByPlatformIDResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	gameInfo, err := db.GetGameInfoByPlatformID(req.PlatformType, req.PlatformID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, GetGameInfoByPlatformIDResponse{
				Status:  "ok",
				Message: "No results found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, GetGameInfoByPlatformIDResponse{
				Status:  "error",
				Message: err.Error(),
			})
		}
	} else {
		c.JSON(http.StatusOK, GetGameInfoByPlatformIDResponse{
			Status:   "ok",
			GameInfo: gameInfo,
		})
	}
}
