package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	groupDeployment = "deployment"
	groupUtility    = "utility"
)

var rootCmd = &cobra.Command{
	Version: "0.1.0",
	Use:     "nessus-on-demand",
	Short:   "Manage just-in-time Nessus deployments in the cloud",
	Long:    `TBD`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddGroup(
		&cobra.Group{ID: groupDeployment, Title: "Deployment Commands"},
		&cobra.Group{ID: groupUtility, Title: "Utility Commands"},
	)

	rootCmd.SetHelpCommandGroupID(groupUtility)
	rootCmd.SetCompletionCommandGroupID(groupUtility)
}
