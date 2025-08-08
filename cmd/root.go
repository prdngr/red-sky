package cmd

import (
	"os"

	"github.com/prdngr/red-sky/core"
	"github.com/spf13/cobra"
)

const (
	groupMain    = "main"
	groupUtility = "utility"
)

var rootCmd = &cobra.Command{
	Use:   "red-sky",
	Short: "Manage just-in-time Nessus deployments in the cloud",
	Long: `RedSky is a handy CLI utility for managing Nessus instances in AWS.
Built using Terraform, RedSky safely bootstraps scanning infrastructure on a per-need basis.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(core.PrintBanner, core.InitNodDir)

	rootCmd.AddGroup(
		&cobra.Group{ID: groupMain, Title: "Main Commands"},
		&cobra.Group{ID: groupUtility, Title: "Utility Commands"},
	)

	rootCmd.SetHelpCommandGroupID(groupUtility)
	rootCmd.SetCompletionCommandGroupID(groupUtility)
}
