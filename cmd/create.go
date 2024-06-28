package cmd

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/prodingerd/nessus-on-demand/core"
	"github.com/spf13/cobra"
)

var (
	region           string
	allowedIp        net.IP
	defaultAllowedIp = net.ParseIP("127.0.0.1")
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a deployment",
	Run:   runCreate,
}

func runCreate(cmd *cobra.Command, args []string) {
	nodDir := core.GetNodDir()
	deploymentId := uuid.New().String()
	tf := core.GetTerraformInstance()

	if err := tf.WorkspaceNew(context.Background(), deploymentId); err != nil {
		log.Fatalf("error creating Terraform workspace: %s", err)
	}

	var variables = []tfexec.PlanOption{
		tfexec.Var("aws_region=" + region),
		tfexec.Var("key_directory=" + nodDir),
		tfexec.Var("deployment_name=" + deploymentId),
	}

	if allowedIp.To4() != nil && !allowedIp.IsLoopback() {
		variables = append(variables, tfexec.Var("allowed_ip="+allowedIp.String()))
	}

	for _, variable := range variables {
		fmt.Println(variable)
	}

	var _, err = tf.Plan(context.Background(), variables...)
	if err != nil {
		log.Fatalf("error planning: %s", err)
	}
	// tf.Apply(context.Background())
}

func init() {
	deploymentCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&region, "region", "r", "eu-central-1", "AWS region to deploy in")
	createCmd.Flags().IPVarP(&allowedIp, "allowed-ip", "a", defaultAllowedIp, `allow-listed IP address (supported "auto", <ipv4_address>)`)
}
