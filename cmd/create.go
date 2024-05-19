package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "create",
	GroupID: groupDeployment,
	Short:   "Creates a deployment",
	Long:    `TBD`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("region", "r", "eu-central-1", "The AWS region to deploy in")
	createCmd.Flags().StringP("allowed-ip", "a", "", `If specified, allow-lists the IP address (supported "auto", "<ipv4_address>")`)
}
