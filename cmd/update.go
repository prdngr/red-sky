package cmd

import (
	"fmt"

	"github.com/prdngr/red-sky/internal"
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update DEPLOYMENT",
	Short: "Update deployment",
	Args:  cobra.ExactArgs(1),
	Run:   runUpdate,
}

func runUpdate(cmd *cobra.Command, args []string) {
	tf := (*internal.Terraform).New(nil)
	tf.UpdateDeployment(args[0], ingressRules)

	internal.PrintHeader("Deployment Summary")

	fmt.Printf("▶ Updated deployment: %s\n", args[0])
}

func init() {
	UpdateCmd.Flags().Var(newIngressRuleSliceValue(nil, &ingressRules), "ingress-rules", "comma-separated list of ingress rules (CIDR:port)")

	UpdateCmd.MarkFlagRequired("ingress-rules")
}
