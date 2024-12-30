package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/fatih/color"
	"github.com/prdngr/nessus-on-demand/core"
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
	details := tf.GetDeploymentDetails()

	core.PrintHeader("Deployment Summary")

	fmt.Printf("Deployment ID: %s\n", details.DeploymentId)
	fmt.Printf("Nessus Interface URL: %s\n", "https://"+details.InstanceIp+":8834")
	fmt.Printf("Allowed IP for Nessus Interface: %s\n", allowedIp)

	core.PrintHeader("Next Steps")

	if allowedIp.IsLoopback() {
		fmt.Print("▶ Forward the Nessus web interface port to your machine using the command below. Then access it via https://localhost:8834.\n")
		color.Cyan("  $ ssh -N -L 8834:127.0.0.1:8834 -i '" + details.SshKeyFile + "' ec2-user@" + details.InstanceIp)
	}

	fmt.Print("▶ Open the Nessus interface in your browser, sign up, and activate your license.\n")
}

func init() {
	deploymentCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&region, "region", "r", region, "AWS region to deploy in")
	createCmd.Flags().IPVar(&allowedIp, "allowed-ip", allowedIp, "allow-listed IP address")
	createCmd.Flags().BoolVar(&autoIp, "auto-ip", autoIp, "automatically determine allow-listed IP")

	createCmd.MarkFlagsMutuallyExclusive("allowed-ip", "auto-ip")
}
