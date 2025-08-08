package cmd

import (
	"fmt"

	"github.com/prdngr/red-sky/internal"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var destroyCmd = &cobra.Command{
	Use:     "destroy DEPLOYMENT [DEPLOYMENT...]",
	Short:   "Destroy deployment(s)",
	Args:    cobra.MinimumNArgs(1),
	GroupID: groupMain,
	Run:     runDestroy,
}

func runDestroy(cmd *cobra.Command, args []string) {
	tf := (*internal.Terraform).New(nil)
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
