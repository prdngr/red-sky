package cmd

import (
	"github.com/spf13/cobra"
)

var deploymentCmd = &cobra.Command{
	Use:     "deployment",
	Short:   "Manage Nessus deployments",
	GroupID: groupMain,
}

func init() {
	rootCmd.AddCommand(deploymentCmd)
}
