package cmd

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/prodingerd/nessus-on-demand/internal"
	"github.com/prodingerd/nessus-on-demand/static"
	"github.com/spf13/cobra"
)

const (
	groupDeployment = "deployment"
	groupUtility    = "utility"
)

var rootCmd = &cobra.Command{
	Version:           internal.NOD_VERSION,
	Use:               "nessus-on-demand",
	Short:             "Manage just-in-time Nessus deployments in the cloud",
	Long:              `TBD`,
	PersistentPreRunE: ensureNodDirectory,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func ensureNodDirectory(cmd *cobra.Command, args []string) error {
	nodDirectory := internal.GetNodDirectory()

	if exists, err := internal.DirectoryExists(nodDirectory); err != nil {
		return err
	} else if exists {
		return nil
	}

	return fs.WalkDir(static.Terraform, "terraform", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		outputPath := filepath.Join(nodDirectory, path)

		if entry.IsDir() {

			if err := internal.CreateDirectoryIfNotExists(outputPath); err != nil {
				return err
			}

			return nil
		}

		if data, err := static.Terraform.ReadFile(path); err != nil {
			return err
		} else {
			if err := os.WriteFile(outputPath, data, os.ModePerm); err != nil {
				return err
			}
		}

		return nil
	})
}

func init() {
	rootCmd.AddGroup(
		&cobra.Group{ID: groupDeployment, Title: "Deployment Commands"},
		&cobra.Group{ID: groupUtility, Title: "Utility Commands"},
	)

	rootCmd.SetHelpCommandGroupID(groupUtility)
	rootCmd.SetCompletionCommandGroupID(groupUtility)
}
