package cmd

import (
	"github.com/nitezs/pcgamedb/log"
	"github.com/nitezs/pcgamedb/task"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type taskCommandConfig struct {
	Crawl bool
}

var taskCmdCfg taskCommandConfig

var taskCmd = &cobra.Command{
	Use:  "task",
	Long: "Start task",
	Run: func(cmd *cobra.Command, args []string) {
		if taskCmdCfg.Crawl {
			task.Crawl(log.Logger)
			c := cron.New()
			_, err := c.AddFunc("0 0 * * *", func() { task.Crawl(log.Logger) })
			if err != nil {
				log.Logger.Error("Failed to add task", zap.Error(err))
			}
			c.Start()
			select {}
		}
	},
}

func init() {
	taskCmd.Flags().BoolVarP(&taskCmdCfg.Crawl, "crawl", "c", false, "enable auto crawl")
	RootCmd.AddCommand(taskCmd)
}
