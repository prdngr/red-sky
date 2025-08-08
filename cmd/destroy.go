package cmd

import (
	"fmt"

	"github.com/prdngr/red-sky/core"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var destroyCmd = &cobra.Command{
	Use:     "destroy [flags] DEPLOYMENT [DEPLOYMENT...]",
	Short:   "Destroy a deployment",
	Args:    cobra.MinimumNArgs(1),
	GroupID: groupMain,
	Run:     runDestroy,
}

func runDestroy(cmd *cobra.Command, args []string) {
	tf := (*core.Terraform).New(nil)
	workspaces := tf.GetWorkspaces()

	for _, deploymentId := range args {
		if !slices.Contains(workspaces, deploymentId) {
			fmt.Printf("Could not find deployment '%s', skipping\n", deploymentId)
			continue
		}

		tf.DestroyDeployment(profile, deploymentId)
	}
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}
