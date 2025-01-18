package cmd

import (
	"fmt"

	"github.com/prdngr/nessus-on-demand/internal"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Nessus on Demand version",
	GroupID: groupUtility,
	Run: func(cmd *cobra.Command, args []string) {
		versionInformation := internal.GetVersion()
		fmt.Printf("Nessus on Demand v%s (%s)\n", versionInformation.BuildVersion, versionInformation.BuildCommit)
		fmt.Printf("Platform: %s\n", versionInformation.Platform)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
