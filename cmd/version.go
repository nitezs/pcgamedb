package cmd

import (
	"fmt"
	"runtime"

	"github.com/nitezs/pcgamedb/constant"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Long:  "Get version of pcgamedb",
	Short: "Get version of pcgamedb",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", constant.Version)
		fmt.Printf("Go: %s\n", runtime.Version())
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
