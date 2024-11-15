package handler

import (
	"net/http"

	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/model"

	"github.com/gin-gonic/gin"
)

type GetGameInfosByNameRequest struct {
	Name string `uri:"name" binding:"required"`
}

type GetGameInfosByNameResponse struct {
	Status    string            `json:"status"`
	Message   string            `json:"message,omitempty"`
	GameInfos []*model.GameInfo `json:"game_infos,omitempty"`
}

// GetGameInfosByName retrieves game information by game name.
// @Summary Retrieve game info by name
// @Description Retrieves game information details by game name
// @Tags game
// @Accept json
// @Produce json
// @Param name path string true "Game Name"
// @Success 200 {object} GetGameInfosByNameResponse
// @Failure 400 {object} GetGameInfosByNameResponse
// @Failure 500 {object} GetGameInfosByNameResponse
// @Router /game/name/{name} [get]
func GetGameInfosByNameHandler(c *gin.Context) {
	var req GetGameInfosByNameRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, GetGameInfosByNameResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	games, err := db.GetGameInfosByName(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetGameInfosByNameResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	if len(games) == 0 {
		c.JSON(http.StatusOK, GetGameInfosByNameResponse{
			Status:  "ok",
			Message: "No results found",
		})
		return
	}
	c.JSON(http.StatusOK, GetGameInfosByNameResponse{
		Status:    "ok",
		GameInfos: games,
	})
}
