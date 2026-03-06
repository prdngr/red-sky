package main

import (
	"os"

	"github.com/prdngr/red-sky/cmd"
	"github.com/prdngr/red-sky/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "red-sky",
	Short:   "The Calm Before the Breach",
	Long:    "RedSky is a handy CLI utility for managing just-in-time offensive infrastructure in AWS.",
	Version: internal.GetVersion().BuildVersion,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(
		internal.ConfigureLogger,
		internal.PrintBanner,
		internal.InitRedSkyDir,
	)

	rootCmd.AddCommand(
		cmd.CreateCmd,
		cmd.DestroyCmd,
		cmd.UpdateCmd,
		cmd.ListCmd,
	)
}
