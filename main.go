package main

import (
	"strings"

	"github.com/nitezs/pcgamedb/cmd"
	"github.com/nitezs/pcgamedb/log"

	"go.uber.org/zap"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		if !strings.Contains(err.Error(), "unknown command") {
			log.Logger.Error("Failed to execute command", zap.Error(err))
		}
	}
}
