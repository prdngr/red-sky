package cmd

import (
	"fmt"

	"github.com/prodingerd/nessus-on-demand/core"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy [flags] DEPLOYMENT [DEPLOYMENT...]",
	Short: "Destroy a deployment",
	Args:  cobra.MinimumNArgs(1),
	Run:   runDestroy,
}

func runDestroy(cmd *cobra.Command, args []string) {
	tf := (*core.Terraform).New(nil)
	workspaces := tf.GetWorkspaces()

	for _, deploymentId := range args {
		if !slices.Contains(workspaces, deploymentId) {
			fmt.Println("Could not find deployment '" + deploymentId + "', skipping")
			continue
		}

		tf.DestroyDeployment(deploymentId)
		tf.DeleteWorkspace(deploymentId)
	}
}

func init() {
	deploymentCmd.AddCommand(destroyCmd)
}
