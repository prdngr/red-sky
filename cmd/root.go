package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/prodingerd/nessus-on-demand/core"
	"github.com/prodingerd/nessus-on-demand/static"

	"github.com/spf13/cobra"
)

const (
	groupMain    = "main"
	groupUtility = "utility"
)

var rootCmd = &cobra.Command{
	Version: core.NodVersion,
	Use:     "nessus-on-demand",
	Short:   "Manage just-in-time Nessus deployments in the cloud",
	Long: `Nessus on Demand is a powerful CLI utility for managing Nessus instances in AWS.
Built using Terraform, Nessus on Demand bootstraps infrastructure with surgical precision.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initNodDirectory, initConfig, initTerraform)

	rootCmd.AddGroup(
		&cobra.Group{ID: groupMain, Title: "Main Commands"},
		&cobra.Group{ID: groupUtility, Title: "Utility Commands"},
	)

	rootCmd.SetHelpCommandGroupID(groupUtility)
	rootCmd.SetCompletionCommandGroupID(groupUtility)
}

func initNodDirectory() {
	nodDirectory := core.GetNodDirectory()

	if exists, err := core.DirectoryExists(nodDirectory); err != nil {
		log.Fatalf("error getting NOD directory: %s", err)
	} else if exists {
		return
	}

	err := fs.WalkDir(static.Embeds, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		outputPath := filepath.Join(nodDirectory, path)

		if entry.IsDir() {
			if err := core.CreateDirectoryIfNotExists(outputPath); err != nil {
				return err
			}

			return nil
		}

		if data, err := static.Embeds.ReadFile(path); err != nil {
			return err
		} else {
			if err := os.WriteFile(outputPath, data, os.ModePerm); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("error initializing NOD directory: %s", err)
	}
}

func initConfig() {
	// viper.AddConfigPath(core.GetConfigDirectory())
	// viper.SetConfigName(core.ConfigName)
	// viper.SetConfigType(core.ConfigType)

	// if err := viper.ReadInConfig(); err != nil {
	// 	log.Fatalf("error reading config file: %s", err)
	// }

	// if err := viper.Unmarshal(&core.Config); err != nil {
	// 	log.Fatalf("error parsing config file: %s", err)
	// }

	// core.K.Unmarshal("", &core.Config)
	core.ReadConfig()
	core.Config.Terraform.Initialized = false
	core.WriteConfig()
	fmt.Println(core.Config.Terraform.Initialized)
}

func initTerraform() {
	if !core.Config.Terraform.Initialized {
		core.InstallTerraform()
		core.InitializeTerraform()

		// viper.Set("terraform.initialized", true)
		// viper.Set("terraform.updatedAt", time.Now().UTC())
		// viper.WriteConfig()
	}
}
