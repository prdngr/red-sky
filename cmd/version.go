package cmd

import (
	"fmt"

	"github.com/prdngr/red-sky/internal"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Show RedSky version",
	GroupID: groupUtility,
	Run: func(cmd *cobra.Command, args []string) {
		versionInformation := internal.GetVersion()
		fmt.Printf("RedSky %s (%s)\n", versionInformation.BuildVersion, versionInformation.BuildCommit)
		fmt.Printf("Platform: %s\n", versionInformation.Platform)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
