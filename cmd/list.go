package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/prodingerd/nessus-on-demand/core"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all deployments",
	Run:   runList,
}

func runList(cmd *cobra.Command, args []string) {
	core.StartSpinner("Initializing NoD")
	tf := core.GetTerraformInstance()
	core.StopSpinner("NoD initialized")

	core.StartSpinner("Retrieving deployments")

	if workspaces, _, err := tf.WorkspaceList(context.Background()); err != nil {
		log.Fatalf("error listing Terraform workspaces: %s", err)
	} else {
		core.StopSpinner("Deployments retrieved")

		if len(workspaces) == 1 {
			fmt.Println("No active deployments")
		}

		for _, workspace := range workspaces {
			if workspace == "default" {
				continue
			}

			fmt.Println(workspace)
		}
	}
}

func init() {
	deploymentCmd.AddCommand(listCmd)
}
