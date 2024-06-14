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
	Short: "lists all deployments",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runList,
}

func runList(cmd *cobra.Command, args []string) {
	core.StartSpinner()

	tf := core.GetTerraformInstance()

	if workspaces, _, err := tf.WorkspaceList(context.Background()); err != nil {
		log.Fatalf("error listing Terraform workspaces: %s", err)
	} else {
		core.StopSpinner()

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
