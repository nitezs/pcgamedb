package handler

import (
	"net/http"
	"pcgamedb/crawler"
	"pcgamedb/model"

	"github.com/gin-gonic/gin"
)

type GetRankingResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message,omitempty"`
	Games   []*model.GameInfo `json:"games"`
}

// GetRanking retrieves game rankings.
// @Summary Retrieve rankings
// @Description Retrieves rankings based on a specified type
// @Tags ranking
// @Accept json
// @Produce json
// @Param type path string true "Ranking Type(top, week-top, best-of-the-year, most-played)"
// @Success 200 {object} GetRankingResponse
// @Failure 400 {object} GetRankingResponse
// @Failure 500 {object} GetRankingResponse
// @Router /ranking/{type} [get]
func GetRankingHandler(c *gin.Context) {
	rankingType, exist := c.Params.Get("type")
	if !exist {
		c.JSON(http.StatusBadRequest, GetRankingResponse{
			Status:  "error",
			Message: "Missing ranking type",
		})
	}
	var f func() ([]*model.GameInfo, error)
	switch rankingType {
	case "top":
		f = crawler.GetSteam250Top250Cache
	case "week-top":
		f = crawler.GetSteam250WeekTop50Cache
	case "best-of-the-year":
		f = crawler.GetSteam250BestOfTheYearCache
	case "most-played":
		f = crawler.GetSteam250MostPlayedCache
	default:
		c.JSON(http.StatusBadRequest, GetRankingResponse{
			Status:  "error",
			Message: "Invalid ranking type",
		})
		return
	}
	rank, err := f()
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetRankingResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, GetRankingResponse{
		Status: "ok",
		Games:  rank,
	})
}
