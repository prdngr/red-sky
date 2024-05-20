package cmd

import (
	"context"
	"fmt"
	"log"
	"path"

	"github.com/hashicorp/go-version"
	install "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hc-install/src"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/prodingerd/nessus-on-demand/internal"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Creates a deployment",
	Long:    `TBD`,
	Run:     runCreate,
	GroupID: groupMain,
}

func runCreate(cmd *cobra.Command, args []string) {
	deploymentId := initDeploymentDirectory()
	fmt.Printf("Deployment ID: %s\n", deploymentId)

	installer := install.NewInstaller()
	terraformVersion := version.MustConstraints(version.NewConstraint("~> 1.8"))
	workingDirectory := path.Join(internal.GetNodDirectory(), internal.NOD_TERRAFORM_DIRECTORY)

	execPath, err := installer.Ensure(context.Background(), []src.Source{
		&fs.Version{
			Product:     product.Terraform,
			ExtraPaths:  []string{workingDirectory},
			Constraints: terraformVersion,
		},
		&releases.LatestVersion{
			Product:     product.Terraform,
			InstallDir:  workingDirectory,
			Constraints: terraformVersion,
		},
	})
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	tf, err := tfexec.NewTerraform(workingDirectory, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	err = tf.WorkspaceNew(context.Background(), deploymentId)
	if err != nil {
		log.Fatalf("error running WorkspaceNew: %s", err)
	}

	// tf.Plan(context.Background())
}

func initDeploymentDirectory() string {
	// deploymentId := uuid.New().String()
	deploymentId := "46d4c9d3-3aab-42d4-bc09-082f5422fcd6"

	return deploymentId
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("region", "r", "eu-central-1", "The AWS region to deploy in")
	createCmd.Flags().StringP("allowed-ip", "a", "none", `Allow-lists an IP address (supported "auto", <ipv4_address>)`)
}
