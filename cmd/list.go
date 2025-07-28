package cmd

import (
	"fmt"

	"github.com/prdngr/nessus-on-demand/core"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all deployments",
	GroupID: groupMain,
	Run:     runList,
}

func runList(cmd *cobra.Command, args []string) {
	tf := (*core.Terraform).New(nil)
	workspaces := tf.GetWorkspaces()

	core.PrintHeader("Deployments")

	if len(workspaces) == 0 {
		fmt.Println("No deployments found")
	}

	for _, workspace := range workspaces {
		fmt.Println(workspace)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}
