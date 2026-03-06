package cmd

import (
	"fmt"

	"github.com/prdngr/red-sky/internal"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var DestroyCmd = &cobra.Command{
	Use:   "destroy DEPLOYMENT [DEPLOYMENT...]",
	Short: "Destroy deployment(s)",
	Args:  cobra.MinimumNArgs(1),
	Run:   runDestroy,
}

func runDestroy(cmd *cobra.Command, args []string) {
	tf := (*internal.Terraform).New(nil)
	workspaces := tf.GetWorkspaces()

	internal.PrintHeader("Deployments")

	for _, deploymentId := range args {
		if !slices.Contains(workspaces, deploymentId) {
			fmt.Printf("▶ Skipped unknown deployment: %s\n", deploymentId)
			continue
		}

		tf.DestroyDeployment(deploymentId)
		fmt.Printf("▶ Destroyed deployment: %s\n", deploymentId)
	}
}
