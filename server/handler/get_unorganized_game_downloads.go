package handler

import (
	"net/http"
	"pcgamedb/db"
	"pcgamedb/model"

	"github.com/gin-gonic/gin"
)

type GetUnorganizedGameDownloadsRequest struct {
	Num int `json:"num" form:"num"`
}

type GetUnorganizedGameDownloadsResponse struct {
	Status        string                `json:"status"`
	Message       string                `json:"message,omitempty"`
	Size          int                   `json:"size,omitempty"`
	GameDownloads []*model.GameDownload `json:"game_downloads,omitempty"`
}

// GetUnorganizedGameDownloads retrieves a list of unorganized game downloads.
// @Summary List unorganized game downloads
// @Description Retrieves game downloads that have not been organized
// @Tags game
// @Accept json
// @Produce json
// @Param num query int false "Number of game downloads to retrieve"
// @Success 200 {object} GetUnorganizedGameDownloadsResponse
// @Failure 400 {object} GetUnorganizedGameDownloadsResponse
// @Failure 500 {object} GetUnorganizedGameDownloadsResponse
// @Router /game/raw/unorganized [get]
func GetUnorganizedGameDownloadsHandler(c *gin.Context) {
	var req GetUnorganizedGameDownloadsRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, GetUnorganizedGameDownloadsResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	if req.Num == 0 || req.Num < 0 {
		req.Num = -1
	}
	gameDownloads, err := db.GetUnorganizedGameDownloads(req.Num)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetUnorganizedGameDownloadsResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	if len(gameDownloads) == 0 {
		c.JSON(http.StatusOK, GetUnorganizedGameDownloadsResponse{
			Status:  "ok",
			Message: "No unorganized game downloads found",
		})
		return
	}
	c.JSON(http.StatusOK, GetUnorganizedGameDownloadsResponse{
		Status:        "ok",
		GameDownloads: gameDownloads,
		Size:          len(gameDownloads),
	})
}
