package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/prodingerd/nessus-on-demand/core"
	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy [flags] DEPLOYMENT [DEPLOYMENT...]",
	Short: "Destroy a deployment",
	Args:  cobra.MinimumNArgs(1),
	Run:   runDestroy,
}

func runDestroy(cmd *cobra.Command, args []string) {
	core.StartSpinner("Initializing NoD")
	tf := core.GetTerraformInstance()
	core.StopSpinner("NoD initialized")

	for _, deploymentId := range args {
		if err := tf.WorkspaceSelect(context.Background(), deploymentId); err != nil {
			fmt.Println("Deployment '" + deploymentId + "' could not be found, skipping")
			continue
		}

		var options = []tfexec.DestroyOption{
			tfexec.Var("aws_region=" + ""),
			tfexec.Var("key_directory=" + ""),
			tfexec.Var("deployment_name=" + ""),
			tfexec.Refresh(false),
		}

		if err := tf.Destroy(context.Background(), options...); err != nil {
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
