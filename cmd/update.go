package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/prdngr/red-sky/internal"
	"github.com/spf13/cobra"
)

var (
	cidrs []net.IPNet
	ports []uint
)

var updateCmd = &cobra.Command{
	Use:     "update DEPLOYMENT",
	Short:   "Update deployment",
	Args:    cobra.ExactArgs(1),
	GroupID: groupMain,
	Run:     runUpdate,
}

func runUpdate(cmd *cobra.Command, args []string) {
	if len(cidrs) != len(ports) {
		log.Fatalf("CIDRs and ports must be of same length.")
	}

	tf := (*internal.Terraform).New(nil)
	tf.UpdateDeployment(args[0], cidrs, ports)

	internal.PrintHeader("Deployment Summary")

	fmt.Printf("▶ Updated deployment: %s\n", args[0])
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().IPNetSliceVar(&cidrs, "cidrs", nil, "8.8.8.8/32[,...]")
	updateCmd.Flags().UintSliceVar(&ports, "ports", nil, "443[,...]")

	updateCmd.MarkFlagRequired("cidrs")
	updateCmd.MarkFlagRequired("ports")
}
