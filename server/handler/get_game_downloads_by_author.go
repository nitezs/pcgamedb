package handler

import (
	"net/http"

	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/model"

	"github.com/gin-gonic/gin"
)

type GetGameDownloadsByAuthorRequest struct {
	Author   string `uri:"author" binding:"required"`
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
}

type GetGameDownloadsByAuthorResponse struct {
	Status        string                `json:"status"`
	Message       string                `json:"message,omitempty"`
	TotalPage     int                   `json:"total_page"`
	GameDownloads []*model.GameDownload `json:"game_downloads,omitempty"`
}

// GetGameDownloadsByAuthorHandler returns all game downloads by author
// @Summary Get game downloads by author
// @Description Get game downloads by author
// @Tags game
// @Accept json
// @Produce json
// @Param author path string true "Author"
// @Param page query int false "Page"
// @Param page_size query int false "Page Size"
// @Success 200 {object} GetGameDownloadsByAuthorResponse
// @Failure 400 {object} GetGameDownloadsByAuthorResponse
// @Failure 500 {object} GetGameDownloadsByAuthorResponse
// @Router /game/raw/author/{author} [get]
func GetGameDownloadsByAuthorHandler(ctx *gin.Context) {
	var req GetGameDownloadsByAuthorRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, GetGameDownloadsByAuthorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, GetGameDownloadsByAuthorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	if req.Page == 0 || req.Page < 0 {
		req.Page = 1
	}
	if req.PageSize == 0 || req.PageSize < 0 {
		req.PageSize = 10
	}
	if req.PageSize > 10 {
		req.PageSize = 10
	}
	downloads, totalPage, err := db.GetGameDownloadsByAuthorPagination(req.Author, req.Page, req.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, GetGameDownloadsByAuthorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	if len(downloads) == 0 {
		ctx.JSON(http.StatusOK, GetGameDownloadsByAuthorResponse{
			Status:  "ok",
			Message: "No results found",
		})
		return
	}
	ctx.JSON(http.StatusOK, GetGameDownloadsByAuthorResponse{
		Status:        "ok",
		TotalPage:     totalPage,
		GameDownloads: downloads,
	})
}
