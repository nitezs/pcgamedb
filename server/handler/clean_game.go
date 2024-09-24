package handler

import (
	"net/http"
	"pcgamedb/log"
	"pcgamedb/task"

	"github.com/gin-gonic/gin"
)

func CleanGameHandler(ctx *gin.Context) {
	task.Clean(log.TaskLogger)
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}
