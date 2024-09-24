package cmd

import (
	"pcgamedb/db"
	"pcgamedb/log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Long:  "List game infos by filter",
	Short: "List game infos by filter",
	Run:   listRun,
}

type listCommandConfig struct {
	Unid bool
}

var listCmdCfg listCommandConfig

func init() {
	listCmd.Flags().BoolVarP(&listCmdCfg.Unid, "unorganized", "u", false, "unorganized")
	RootCmd.AddCommand(listCmd)
}

func listRun(cmd *cobra.Command, args []string) {
	if listCmdCfg.Unid {
		games, err := db.GetUnorganizedGameDownloads(-1)
		if err != nil {
			log.Logger.Error("Failed to get games", zap.Error(err))
		}
		for _, game := range games {
			log.Logger.Info(
				"Game",
				zap.Any("game_id", game.ID),
				zap.String("raw_name", game.RawName),
				zap.String("name", game.Name),
			)
		}
	}
}
