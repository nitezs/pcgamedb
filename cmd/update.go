package cmd

import (
	"pcgamedb/crawler"
	"pcgamedb/db"
	"pcgamedb/log"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Long:  "Update game info by game data platform",
	Short: "Update game info by game data platform",
	Run:   updateRun,
}

type updateCommandConfig struct {
	PlatformID int
	Platform   string
	ID         string
}

var updateCmdcfx updateCommandConfig

func init() {
	updateCmd.Flags().IntVarP(&updateCmdcfx.PlatformID, "platform-id", "p", 0, "platform id")
	updateCmd.Flags().StringVarP(&updateCmdcfx.Platform, "platform", "t", "", "platform")
	updateCmd.Flags().StringVarP(&updateCmdcfx.ID, "game-id", "i", "", "game info id")
	RootCmd.AddCommand(updateCmd)
}

func updateRun(cmd *cobra.Command, args []string) {
	id, err := primitive.ObjectIDFromHex(updateCmdcfx.ID)
	if err != nil {
		log.Logger.Error("Failed to parse game info id", zap.Error(err))
		return
	}
	oldInfo, err := db.GetGameInfoByID(id)
	if err != nil {
		log.Logger.Error("Failed to get game info", zap.Error(err))
		return
	}
	newInfo, err := crawler.GenerateGameInfo(updateCmdcfx.Platform, updateCmdcfx.PlatformID)
	if err != nil {
		log.Logger.Error("Failed to generate game info", zap.Error(err))
		return
	}
	newInfo.ID = id
	newInfo.GameIDs = oldInfo.GameIDs
	err = db.SaveGameInfo(newInfo)
	if err != nil {
		log.Logger.Error("Failed to save game info", zap.Error(err))
	}
}
