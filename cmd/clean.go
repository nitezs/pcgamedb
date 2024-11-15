package cmd

import (
	"github.com/nitezs/pcgamedb/log"
	"github.com/nitezs/pcgamedb/task"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Long:  "Clean database",
	Short: "Clean database",
	Run: func(cmd *cobra.Command, args []string) {
		task.Clean(log.Logger)
	},
}

func init() {
	RootCmd.AddCommand(cleanCmd)
}
