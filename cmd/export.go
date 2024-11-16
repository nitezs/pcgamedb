package cmd

import (
	"os"
	"path/filepath"

	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/log"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
)

type exportCommandConfig struct {
	output string
}

var exportCmdCfg exportCommandConfig

var exportCmd = &cobra.Command{
	Use:   "export",
	Long:  "Export data to json files",
	Short: "Export data to json files",
	Run: func(cmd *cobra.Command, args []string) {
		infoJson, gameJson, err := db.Export()
		if err != nil {
			log.Logger.Error("Error exporting data", zap.Error(err))
			return
		}
		infoFilePath := filepath.Join(exportCmdCfg.output, "game_infos.json")
		gameFilePath := filepath.Join(exportCmdCfg.output, "games.json")

		err = os.WriteFile(infoFilePath, infoJson, 0644)
		if err != nil {
			if err != nil {
				log.Logger.Error("Error exporting data", zap.Error(err))
				return
			}
		}
		err = os.WriteFile(gameFilePath, gameJson, 0644)
		if err != nil {
			if err != nil {
				log.Logger.Error("Error exporting data", zap.Error(err))
				return
			}
		}
		log.Logger.Info("Data exported successfully")
	},
}

func init() {
	exportCmd.Flags().StringVarP(&exportCmdCfg.output, "output", "o", "", "Output directory for json files")
	RootCmd.AddCommand(exportCmd)
}
