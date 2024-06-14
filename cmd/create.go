package cmd

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/prodingerd/nessus-on-demand/core"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "creates a deployment",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runCreate,
}

func runCreate(cmd *cobra.Command, args []string) {
	deploymentId := uuid.New().String()
	tf := core.GetTerraformInstance()

	if err := tf.WorkspaceNew(context.Background(), deploymentId); err != nil {
		log.Fatalf("error creating Terraform workspace: %s", err)
	}

	var variables = []tfexec.PlanOption{
		tfexec.Var("aws_region=eu-central-1"),
		tfexec.Var("deployment_name=" + deploymentId),
	}

	var _, _ = tf.Plan(context.Background(), variables...)
	// tf.Apply(context.Background())
}

func init() {
	deploymentCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("region", "r", "eu-central-1", "The AWS region to deploy in")
	createCmd.Flags().StringP("allowed-ip", "a", "none", `Allow-lists an IP address (supported "auto", <ipv4_address>)`)
}
