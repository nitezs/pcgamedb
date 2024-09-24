package handler

import (
	"net/http"
	"pcgamedb/db"

	"github.com/gin-gonic/gin"
)

type GetAllAuthorsResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message,omitempty"`
	Authors []string `json:"authors,omitempty"`
}

// GetAllAuthorsHandler returns all authors
// @Summary Get all authors
// @Description Get all authors
// @Tags author
// @Accept json
// @Produce json
// @Success 200 {object} GetAllAuthorsResponse
// @Failure 500 {object} GetAllAuthorsResponse
// @Router /author [get]
func GetAllAuthorsHandler(ctx *gin.Context) {
	authors, err := db.GetAllAuthors()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, GetAllAuthorsResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	if len(authors) == 0 {
		ctx.JSON(http.StatusOK, GetAllAuthorsResponse{
			Status:  "ok",
			Message: "No authors found",
		})
		return
	}
	ctx.JSON(http.StatusOK, GetAllAuthorsResponse{
		Status:  "ok",
		Authors: authors,
	})
}
