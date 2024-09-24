package main

import (
	"pcgamedb/cmd"
	"pcgamedb/log"
	"strings"

	"go.uber.org/zap"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		if !strings.Contains(err.Error(), "unknown command") {
			log.Logger.Error("Failed to execute command", zap.Error(err))
		}
	}
}
