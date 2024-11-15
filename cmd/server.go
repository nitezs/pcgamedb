package cmd

import (
	"github.com/nitezs/pcgamedb/config"
	"github.com/nitezs/pcgamedb/server"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Long:  "Start API server",
	Short: "Start API server",
	Run:   ServerRun,
}

type serverCommandConfig struct {
	Port      string
	AutoCrawl bool
}

var serverCmdCfg serverCommandConfig

func init() {
	serverCmd.Flags().StringVarP(&serverCmdCfg.Port, "port", "p", "8080", "server port")
	serverCmd.Flags().BoolVarP(&serverCmdCfg.AutoCrawl, "auto-crawl", "c", true, "enable auto crawl")
	RootCmd.AddCommand(serverCmd)
}

func ServerRun(cmd *cobra.Command, args []string) {
	if serverCmdCfg.AutoCrawl {
		config.Config.Server.AutoCrawl = true
	}
	config.Config.Server.Port = serverCmdCfg.Port
	server.Run()
}
