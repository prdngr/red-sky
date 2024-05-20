package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use:     "destroy DEPLOYMENT",
	Short:   "Destroys a deployment",
	Long:    `TBD`,
	Args:    cobra.MinimumNArgs(1),
	Run:     runDestroy,
	GroupID: groupMain,
}

func runDestroy(cmd *cobra.Command, args []string) {
	fmt.Println("destroy called")
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}
