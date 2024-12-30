package cmd

import (
	"log"
	"net"

	"github.com/prodingerd/nessus-on-demand/core"
	"github.com/spf13/cobra"
)

var (
	region    = "eu-central-1"
	allowedIp = net.ParseIP("127.0.0.1")
	autoIp    = false
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a deployment",
	Run:   runCreate,
}

func runCreate(cmd *cobra.Command, args []string) {
	if autoIp {
		if publicIp, err := core.GetPublicIp(); err != nil {
			log.Fatalf("error determining public IP address: %s", err)
		} else {
			allowedIp = publicIp
		}
	}

	tf := (*core.Terraform).New(nil)
	workspace := tf.CreateWorkspace()

	tf.ApplyDeployment(workspace, region, allowedIp)
	tf.GetOutput()
}

func init() {
	deploymentCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&region, "region", "r", region, "AWS region to deploy in")
	createCmd.Flags().IPVar(&allowedIp, "allowed-ip", allowedIp, "allow-listed IP address")
	createCmd.Flags().BoolVar(&autoIp, "auto-ip", autoIp, "automatically determine allow-listed IP")

	createCmd.MarkFlagsMutuallyExclusive("allowed-ip", "auto-ip")
}
