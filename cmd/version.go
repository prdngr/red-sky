package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Nessus on Demand version",
	GroupID: groupUtility,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("nessus-on-demand version " + version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
