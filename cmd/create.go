package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/fatih/color"
	"github.com/prdngr/red-sky/internal"
	"github.com/spf13/cobra"
)

var (
	profile        = "default"
	region         = "eu-central-1"
	autoIp         = false
	allowedIp      = net.ParseIP("127.0.0.1")
	deploymentType string
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a deployment",
	GroupID: groupMain,
	Run:     runCreate,
}

func runCreate(cmd *cobra.Command, args []string) {
	// TODO Improve input validation and error handling.
	if deploymentType != "nessus" && deploymentType != "kali" && deploymentType != "c2" {
		log.Fatalf("Invalid deployment type: '%s'. Allowed types are 'nessus', 'kali', or 'c2'.", deploymentType)
	}

	internal.InitializeAwsSession(profile)

	if autoIp {
		if publicIp, err := internal.GetPublicIp(); err != nil {
			log.Fatalf("error determining public IP address: %s", err)
		} else {
			allowedIp = publicIp
		}
	}

	tf := (*internal.Terraform).New(nil)

	tf.ApplyDeployment(profile, region, deploymentType, allowedIp)
	details := tf.GetDeploymentDetails()

	internal.PrintHeader("Deployment Summary")

	fmt.Printf("Deployment ID: %s\n", details.DeploymentId)
	fmt.Printf("Nessus Interface: %s\n", "https://"+details.InstanceIp+":8834")
	fmt.Printf("Allowed IP Address: %s\n", allowedIp)

	internal.PrintHeader("Next Steps")

	fmt.Println("▶ Connect to the instance via SSH using the command below.")
	color.Cyan("  $ ssh -i '%s' ec2-user@%s", details.SshKeyFile, details.InstanceIp)

	if deploymentType == "nessus" {
		if allowedIp.IsLoopback() {
			fmt.Println("▶ Forward the Nessus web interface port to your machine using the command below. Then access it via https://localhost:8834.")
			color.Cyan("  $ ssh -N -L 8834:127.0.0.1:8834 -i '%s' ec2-user@%s", details.SshKeyFile, details.InstanceIp)
		}
		fmt.Println("▶ Open the Nessus interface in your browser, sign up, and activate your license.")
	}
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&profile, "profile", "p", profile, "AWS profile")
	createCmd.Flags().StringVarP(&region, "region", "r", region, "AWS region")
	createCmd.Flags().StringVarP(&deploymentType, "type", "t", "", "Deployment type (nessus, kali, or c2)")
	createCmd.Flags().IPVar(&allowedIp, "allowed-ip", allowedIp, "allow-listed IP address")
	createCmd.Flags().BoolVar(&autoIp, "auto-ip", autoIp, "automatically determine allow-listed IP")

	createCmd.MarkFlagRequired("type")
	createCmd.MarkFlagsMutuallyExclusive("allowed-ip", "auto-ip")
}
