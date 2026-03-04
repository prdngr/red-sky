package cmd

import (
	"os"

	"github.com/prdngr/red-sky/internal"
	"github.com/spf13/cobra"
)

const (
	groupMain    = "main"
	groupUtility = "utility"
)

var rootCmd = &cobra.Command{
	Use:   "red-sky",
	Short: "The Calm Before the Breach",
	Long:  "RedSky is a handy CLI utility for managing just-in-time offensive infrastructure in AWS.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(
		internal.ConfigureLogger,
		internal.PrintBanner,
		internal.InitRedSkyDir,
	)

	rootCmd.AddGroup(
		&cobra.Group{ID: groupMain, Title: "Main Commands"},
		&cobra.Group{ID: groupUtility, Title: "Utility Commands"},
	)

	rootCmd.SetHelpCommandGroupID(groupUtility)
	rootCmd.SetCompletionCommandGroupID(groupUtility)
}
