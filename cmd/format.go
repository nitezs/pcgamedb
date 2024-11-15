package cmd

import (
	"strings"

	"github.com/nitezs/pcgamedb/crawler"
	"github.com/nitezs/pcgamedb/db"
	"github.com/nitezs/pcgamedb/log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Format game downloads name by formatter",
	Long:  "Format game downloads name by formatter",
	Run:   formatRun,
}

type FormatCommandConfig struct {
	Source string
}

var formatCmdCfg FormatCommandConfig

func init() {
	formatCmd.Flags().StringVarP(&formatCmdCfg.Source, "source", "s", "", "source to fix (fitgirl/dodi/kaoskrew/freegog/xatab/onlinefix/armgddn)")
	RootCmd.AddCommand(formatCmd)
}

func formatRun(cmd *cobra.Command, args []string) {
	formatCmdCfg.Source = strings.ToLower(formatCmdCfg.Source)
	switch formatCmdCfg.Source {
	case "dodi":
		items, err := db.GetDODIGameDownloads()
		if err != nil {
			log.Logger.Error("Failed to get games", zap.Error(err))
			return
		}
		for _, item := range items {
			oldName := item.Name
			item.Name = crawler.DODIFormatter(item.RawName)
			if oldName != item.Name {
				log.Logger.Info("Fix name", zap.String("old", oldName), zap.String("raw", item.RawName), zap.String("name", item.Name))
				err := db.SaveGameDownload(item)
				if err != nil {
					log.Logger.Error("Failed to update item", zap.Error(err))
				}
			}
		}
	case "kaoskrew":
		items, err := db.GetKaOsKrewGameDownloads()
		if err != nil {
			log.Logger.Error("Failed to get games", zap.Error(err))
			return
		}
		for _, item := range items {
			oldName := item.Name
			item.Name = crawler.KaOsKrewFormatter(item.RawName)
			if oldName != item.Name {
				log.Logger.Info("Fix name", zap.String("old", oldName), zap.String("raw", item.RawName), zap.String("name", item.Name))
				err := db.SaveGameDownload(item)
				if err != nil {
					log.Logger.Error("Failed to update item", zap.Error(err))
				}
			}
		}
	case "freegog":
		items, err := db.GetFreeGOGGameDownloads()
		if err != nil {
			log.Logger.Error("Failed to get games", zap.Error(err))
			return
		}
		for _, item := range items {
			oldName := item.Name
			item.Name = crawler.FreeGOGFormatter(item.RawName)
			if oldName != item.Name {
				log.Logger.Info("Fix name", zap.String("old", oldName), zap.String("raw", item.RawName), zap.String("name", item.Name))
				err := db.SaveGameDownload(item)
				if err != nil {
					log.Logger.Error("Failed to update item", zap.Error(err))
				}
			}
		}
	case "xatab":
		items, err := db.GetXatabGameDownloads()
		if err != nil {
			log.Logger.Error("Failed to get games", zap.Error(err))
			return
		}
		for _, item := range items {
			oldName := item.Name
			item.Name = crawler.XatabFormatter(item.RawName)
			if oldName != item.Name {
				log.Logger.Info("Fix name", zap.String("old", oldName), zap.String("raw", item.RawName), zap.String("name", item.Name))
				err := db.SaveGameDownload(item)
				if err != nil {
					log.Logger.Error("Failed to update item", zap.Error(err))
				}
			}
		}
	case "onlinefix":
		items, err := db.GetOnlineFixGameDownloads()
		if err != nil {
			log.Logger.Error("Failed to get games", zap.Error(err))
			return
		}
		for _, item := range items {
			oldName := item.Name
			item.Name = crawler.OnlineFixFormatter(item.RawName)
			if oldName != item.Name {
				log.Logger.Info("Fix name", zap.String("old", oldName), zap.String("raw", item.RawName), zap.String("name", item.Name))
				err := db.SaveGameDownload(item)
				if err != nil {
					log.Logger.Error("Failed to update item", zap.Error(err))
				}
			}
		}
	case "armgddn":
		items, err := db.GetARMGDDNGameDownloads()
		if err != nil {
			log.Logger.Error("Failed to get games", zap.Error(err))
			return
		}
		for _, item := range items {
			oldName := item.Name
			item.Name = crawler.ARMGDDNFormatter(item.RawName)
			if oldName != item.Name {
				log.Logger.Info("Fix name", zap.String("old", oldName), zap.String("raw", item.RawName), zap.String("name", item.Name))
				err := db.SaveGameDownload(item)
				if err != nil {
					log.Logger.Error("Failed to update item", zap.Error(err))
				}
			}
		}
	}
}
