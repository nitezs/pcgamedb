package cmd

import (
	"github.com/nitezs/pcgamedb/crawler"
	"github.com/nitezs/pcgamedb/log"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
)

var supplementCmd = &cobra.Command{
	Use:   "supplement",
	Long:  "Supplement platform id to game info",
	Short: "Supplement platform id to game info",
	Run: func(cmd *cobra.Command, args []string) {
		err := crawler.SupplementPlatformIDToGameInfo(log.Logger)
		if err != nil {
			log.Logger.Error("Error supplementing platform id to game info", zap.Error(err))
		}
	},
}

func init() {
	RootCmd.AddCommand(supplementCmd)
}
