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
	region    string = "eu-central-1"
	allowedIp net.IP = net.ParseIP("127.0.0.1")
	autoIp    bool   = false
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a deployment",
	Run:   runCreate,
}

func runCreate(cmd *cobra.Command, args []string) {
	nodDir := core.GetNodDir()
	deploymentId := uuid.New().String()

	var options = []tfexec.ApplyOption{
		tfexec.Var("aws_region=" + region),
		tfexec.Var("key_directory=" + nodDir),
		tfexec.Var("deployment_name=" + deploymentId),
	}

	if allowedIp.To4() != nil && !allowedIp.IsLoopback() {
		options = append(options, tfexec.Var("allowed_ip="+allowedIp.String()))
	} else if autoIp {
		if publicIp, err := core.GetPublicIp(); err != nil {
			log.Fatalf("error determining allowed IP: %s", err)
		} else {
			options = append(options, tfexec.Var("allowed_ip="+publicIp.String()))
		}
	}

	core.StartSpinner("Initializing NoD")
	tf := core.GetTerraformInstance()
	core.StopSpinner("NoD initialized")

	if err := tf.WorkspaceNew(context.Background(), deploymentId); err != nil {
		log.Fatalf("error creating Terraform workspace: %s", err)
	}

	// TODO Remove debug output.
	for _, variable := range options {
		fmt.Println(variable)
	}

	core.StartSpinner("Deploying Nessus")

	// if _, err := tf.Plan(context.Background(), options...); err != nil {
	// 	core.StopSpinner("Could not plan deployment")

	// 	if err := tf.WorkspaceDelete(context.Background(), deploymentId); err != nil {
	// 		log.Fatalf("error deleting Terraform workspace: %s", err)
	// 	}

	// 	return
	// }

	if tf.Apply(context.Background(), options...) != nil {
		core.StopSpinner("Deployment failed")

		if err := tf.WorkspaceDelete(context.Background(), deploymentId); err != nil {
			log.Fatalf("error deleting Terraform workspace: %s", err)
		}

		return
	}

	core.StopSpinner("Nessus deployed")

	if outputs, err := tf.Output(context.Background()); err != nil {
		log.Fatalf("error retrieving Terraform output: %s", err)
	} else {
		for _, output := range outputs {
			fmt.Println(output.Value)
		}
	}
}

func init() {
	deploymentCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&region, "region", "r", region, "AWS region to deploy in")
	createCmd.Flags().IPVar(&allowedIp, "allowed-ip", allowedIp, "allow-listed IP address")
	createCmd.Flags().BoolVar(&autoIp, "auto-ip", autoIp, "automatically determine allow-listed IP")

	createCmd.MarkFlagsMutuallyExclusive("allowed-ip", "auto-ip")
}
