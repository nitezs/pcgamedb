package cmd

import (
	"encoding/json"
	"os"
	"pcgamedb/crawler"
	"pcgamedb/db"
	"pcgamedb/log"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

var addCmd = &cobra.Command{
	Use:   "manual",
	Long:  "Manually add information for games that cannot match IDs from IGDB, Steam or GOG",
	Short: "Manually add information for games that cannot match IDs from IGDB, Steam or GOG",
	Run:   addRun,
}

type ManualCommandConfig struct {
	GameID     string `json:"game_id"`
	Platform   string `json:"platform"`
	PlatformID int    `json:"platform_id"`
	Config     string
}

var manualCmdCfg ManualCommandConfig

func init() {
	addCmd.Flags().StringVarP(&manualCmdCfg.GameID, "game-id", "i", "", "repack game id")
	addCmd.Flags().StringVarP(&manualCmdCfg.Platform, "platform", "t", "", "platform")
	addCmd.Flags().IntVarP(&manualCmdCfg.PlatformID, "platform-id", "p", 0, "platform id")
	addCmd.Flags().StringVarP(&manualCmdCfg.Config, "config", "c", "", "config path")
	organizeCmd.AddCommand(addCmd)
}

func addRun(cmd *cobra.Command, args []string) {
	c := []*ManualCommandConfig{}
	if manualCmdCfg.Config != "" {
		data, err := os.ReadFile(manualCmdCfg.Config)
		if err != nil {
			log.Logger.Error("Failed to read config file", zap.Error(err))
			return
		}
		if err = json.Unmarshal(data, &c); err != nil {
			log.Logger.Error("Failed to unmarshal config file", zap.Error(err))
			return
		}
	} else {
		c = append(c, &manualCmdCfg)
	}
	for _, v := range c {
		objID, err := primitive.ObjectIDFromHex(v.GameID)
		if err != nil {
			log.Logger.Error("Failed to parse game id", zap.Error(err))
			continue
		}
		info, err := crawler.OrganizeGameDownloadManually(objID, v.Platform, v.PlatformID)
		if err != nil {
			log.Logger.Error("Failed to add game info", zap.Error(err))
			continue
		}
		err = db.SaveGameInfo(info)
		if err != nil {
			log.Logger.Error("Failed to save game info", zap.Error(err))
			continue
		}
		log.Logger.Info("Added game info", zap.String("game_id", v.GameID), zap.String("id_type", v.Platform), zap.Int("id", v.PlatformID))
	}
}
