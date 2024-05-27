package cmd

import (
	"context"
	"log"

	"github.com/prodingerd/nessus-on-demand/core"
	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy [flags] DEPLOYMENT [DEPLOYMENT...]",
	Short: "destroys a deployment",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run:  runDestroy,
}

func runDestroy(cmd *cobra.Command, args []string) {
	tf := core.InitializeTerraform()

	for _, deploymentId := range args {
		if err := tf.WorkspaceSelect(context.Background(), deploymentId); err != nil {
			log.Fatalf("error selecting Terraform workspace")
		}

		if err := tf.Destroy(context.Background()); err != nil {
			log.Fatalf("error destroying Terraform deployment: %s", err)
		}

		if err := tf.WorkspaceDelete(context.Background(), deploymentId); err != nil {
			log.Fatalf("error deleting Terraform workspace: %s", err)
		}
	}
}

func init() {
	deploymentCmd.AddCommand(destroyCmd)
}
