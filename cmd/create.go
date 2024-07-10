package cmd

import (
	"log"
	"net"

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

	var options = []tfexec.ApplyOption{
		tfexec.Var("aws_region=" + region),
		tfexec.Var("key_directory=" + nodDir),
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

	tf := (*core.Terraform).New(nil)
	workspace := tf.CreateWorkspace()

	tf.ApplyDeployment(workspace, append(options, tfexec.Var("deployment_name="+workspace)))
}

func init() {
	deploymentCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&region, "region", "r", region, "AWS region to deploy in")
	createCmd.Flags().IPVar(&allowedIp, "allowed-ip", allowedIp, "allow-listed IP address")
	createCmd.Flags().BoolVar(&autoIp, "auto-ip", autoIp, "automatically determine allow-listed IP")

	createCmd.MarkFlagsMutuallyExclusive("allowed-ip", "auto-ip")
}
