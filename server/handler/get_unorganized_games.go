package handler

import (
	"net/http"

	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/model"

	"github.com/gin-gonic/gin"
)

type GetUnorganizedGameItemsRequest struct {
	Num int `json:"num" form:"num"`
}

type GetUnorganizedGameItemsResponse struct {
	Status    string            `json:"status"`
	Message   string            `json:"message,omitempty"`
	Size      int               `json:"size,omitempty"`
	GameItems []*model.GameItem `json:"game_downloads,omitempty"`
}

// GetUnorganizedGameItems retrieves a list of unorganized game downloads.
// @Summary List unorganized game downloads
// @Description Retrieves game downloads that have not been organized
// @Tags game
// @Accept json
// @Produce json
// @Param num query int false "Number of game downloads to retrieve"
// @Success 200 {object} GetUnorganizedGameItemsResponse
// @Failure 400 {object} GetUnorganizedGameItemsResponse
// @Failure 500 {object} GetUnorganizedGameItemsResponse
// @Router /game/raw/unorganized [get]
func GetUnorganizedGameItemsHandler(c *gin.Context) {
	var req GetUnorganizedGameItemsRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, GetUnorganizedGameItemsResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	if req.Num == 0 || req.Num < 0 {
		req.Num = -1
	}
	gameDownloads, err := db.GetUnorganizedGameItems(req.Num)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetUnorganizedGameItemsResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	if len(gameDownloads) == 0 {
		c.JSON(http.StatusOK, GetUnorganizedGameItemsResponse{
			Status:  "ok",
			Message: "No unorganized game downloads found",
		})
		return
	}
	c.JSON(http.StatusOK, GetUnorganizedGameItemsResponse{
		Status:    "ok",
		GameItems: gameDownloads,
		Size:      len(gameDownloads),
	})
}
