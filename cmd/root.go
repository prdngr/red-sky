package cmd

import (
	"os"

	"github.com/prodingerd/nessus-on-demand/core"
	"github.com/spf13/cobra"
)

const (
	groupMain    = "main"
	groupUtility = "utility"
)

var rootCmd = &cobra.Command{
	Version: "0.1.0",
	Use:     "nessus-on-demand",
	Short:   "Manage just-in-time Nessus deployments in the cloud",
	Long: `Nessus on Demand is a powerful CLI utility for managing Nessus instances in AWS.
Built using Terraform, Nessus on Demand bootstraps infrastructure with surgical precision.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(core.InitNodDirectory)

	rootCmd.AddGroup(
		&cobra.Group{ID: groupMain, Title: "Main Commands"},
		&cobra.Group{ID: groupUtility, Title: "Utility Commands"},
	)

	rootCmd.SetHelpCommandGroupID(groupUtility)
	rootCmd.SetCompletionCommandGroupID(groupUtility)
}
