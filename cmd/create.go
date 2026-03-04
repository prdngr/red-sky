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
	profile        string
	region         = "eu-central-1"
	adminCidr      net.IPNet
	autoAdminCidr  = false
	deploymentType DeploymentType
	ingressRules   []internal.IngressRule
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create deployment",
	GroupID: groupMain,
	Run:     runCreate,
}

func runCreate(cmd *cobra.Command, args []string) {
	internal.InitializeAwsSession(profile, region)

	if autoAdminCidr {
		if _, publicCidr, err := internal.GetPublicIp(); err != nil {
			log.Fatalf("error determining public IP address: %s", err)
		} else {
			adminCidr = *publicCidr
		}
	}

	tf := (*internal.Terraform).New(nil)

	tf.ApplyDeployment(profile, region, deploymentType.String(), adminCidr, ingressRules)
	details := tf.GetDeploymentDetails()

	internal.PrintHeader("Deployment Summary")

	fmt.Printf("▶ Deployment ID: %s\n", details.DeploymentId)
	fmt.Printf("▶ Allowed admin CIDR: %s\n", adminCidr.String())

	internal.PrintHeader("Connection Details")

	switch deploymentType {
	case deploymentTypeNessus:
		if adminCidr.IP == nil {
			fmt.Println("▶ Use the following command to forward the Nessus web interface locally, then access it via https://localhost:8834:")
			color.Cyan("  $ ssh -N -L 8834:127.0.0.1:8834 -i '%s' ec2-user@%s", details.SshKeyFile, details.InstanceIp)
		} else {
			fmt.Printf("▶ Access the Nessus web interface via https://%s:8834\n", details.InstanceIp)
		}

		fmt.Println("▶ Use the following command to SSH into the Nessus instance:")
		color.Cyan("  $ ssh -i '%s' ec2-user@%s", details.SshKeyFile, details.InstanceIp)
	case deploymentTypeKali:
		fmt.Println("▶ Use the following command to SSH into the Kali instance:")
		color.Cyan("  $ ssh -i '%s' kali@%s", details.SshKeyFile, details.InstanceIp)
	case deploymentTypeC2:
		if adminCidr.IP == nil {
			fmt.Println("▶ Use the following command to forward the C2 web interface locally, then access it via https://localhost:8834:")
			color.Cyan("  $ ssh -N -L 7443:127.0.0.1:7443 -i '%s' kali@%s", details.SshKeyFile, details.InstanceIp)
		} else {
			fmt.Printf("▶ Access the C2 web interface via https://%s:7443\n", details.InstanceIp)
		}

		fmt.Printf("▶ Use the following URL for HTTPS callbacks: %s\n", details.CloudFrontUrl)
		fmt.Println("▶ Use the following command to SSH into the C2 instance:")
		color.Cyan("  $ ssh -i '%s' kali@%s", details.SshKeyFile, details.InstanceIp)
	}
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS profile")
	createCmd.Flags().StringVarP(&region, "region", "r", region, "AWS region")
	createCmd.Flags().VarP(&deploymentType, "type", "t", `deployment type ("nessus", "kali", or "c2")`)
	createCmd.Flags().IPNetVar(&adminCidr, "admin-cidr", adminCidr, "allow-listed admin CIDR")
	createCmd.Flags().BoolVar(&autoAdminCidr, "auto-admin-cidr", autoAdminCidr, "auto determine the admin CIDR (true or false)")
	createCmd.Flags().Var(newIngressRuleSliceValue(nil, &ingressRules), "ingress-rules", "additional ingress rules (CIDR:port)")

	createCmd.MarkFlagRequired("type")
	createCmd.MarkFlagsMutuallyExclusive("admin-cidr", "auto-admin-cidr")
}
