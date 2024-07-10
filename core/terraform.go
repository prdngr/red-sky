package core

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/adrg/xdg"
	"github.com/google/uuid"
	"github.com/hashicorp/go-version"
	install "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hc-install/src"
	"github.com/hashicorp/terraform-exec/tfexec"
)

const (
	defaultWorkspace    = "default"
	terraformVersion    = "~> 1.8"
	terraformWorkingDir = NodDir + "terraform/"
	terraformInstallDir = NodDir + "bin/x"
)

type Terraform struct {
	instance *tfexec.Terraform
}

func (*Terraform) New() *Terraform {
	StartSpinner("Initializing NoD")

	tf, err := tfexec.NewTerraform(getTerraformWorkingDir(), getTerraformExecutable())
	if err != nil {
		log.Fatalf("error creating Terraform instance: %s", err)
	}

	if err := tf.Init(context.Background(), tfexec.Backend(false)); err != nil {
		log.Fatalf("error initializing Terraform: %s", err)
	}

	StopSpinner("NoD initialized")

	return &Terraform{instance: tf}
}

func (tf *Terraform) GetWorkspaces() []string {
	StartSpinner("Retrieving deployments")

	workspaces, _, err := tf.instance.WorkspaceList(context.Background())
	if err != nil {
		log.Fatalf("error listing Terraform workspaces: %s", err)
	}

	var nonDefaultWorkspaces []string

	for _, workspace := range workspaces {
		if workspace == defaultWorkspace {
			continue
		}

		nonDefaultWorkspaces = append(nonDefaultWorkspaces, workspace)
	}

	StopSpinner("Deployments retrieved")

	return nonDefaultWorkspaces
}

func (tf *Terraform) CreateWorkspace() string {
	workspace := uuid.New().String()

	if err := tf.instance.WorkspaceNew(context.Background(), workspace); err != nil {
		log.Fatalf("error creating Terraform workspace: %s", err)
	}

	return workspace
}

func (tf *Terraform) DeleteWorkspace(workspace string) {
	if err := tf.instance.WorkspaceSelect(context.Background(), defaultWorkspace); err != nil {
		log.Fatalf("error selecting Terraform workspace: %s", err)
	}

	if err := tf.instance.WorkspaceDelete(context.Background(), workspace); err != nil {
		log.Fatalf("error deleting Terraform workspace: %s", err)
	}
}

func (tf *Terraform) ApplyDeployment(workspace string, region string, allowedIp net.IP) {
	StartSpinner("Deploying Nessus")

	var options = []tfexec.ApplyOption{
		createVar("aws_region", region),
		createVar("key_directory", GetNodDir()),
		createVar("deployment_name", workspace),
	}

	if allowedIp.To4() != nil && !allowedIp.IsLoopback() {
		options = append(options, createVar("allowed_ip", allowedIp.String()))
	}

	if tf.instance.Apply(context.Background(), options...) != nil {
		StopSpinnerError("Deployment failed")
		tf.DeleteWorkspace(workspace)
		return
	}

	if outputs, err := tf.instance.Output(context.Background()); err != nil {
		log.Fatalf("error retrieving Terraform output: %s", err)
	} else {
		for _, output := range outputs {
			fmt.Println(output.Value)
		}
	}

	StopSpinner("Nessus deployed")
}

func (tf *Terraform) DestroyDeployment(workspace string) {
	if err := tf.instance.WorkspaceSelect(context.Background(), workspace); err != nil {
		log.Fatalf("error selecting Terraform workspace: %s", err)
	}

	var options = []tfexec.DestroyOption{
		createVar("aws_region", ""),
		createVar("key_directory", GetNodDir()),
		createVar("deployment_name", workspace),
		// tfexec.Refresh(false),
	}

	if err := tf.instance.Destroy(context.Background(), options...); err != nil {
		log.Fatalf("error destroying Terraform deployment: %s", err)
	}
}

func createVar(key string, value string) *tfexec.VarOption {
	return tfexec.Var(key + "=" + value)
}

func getTerraformWorkingDir() string {
	workingDir, err := xdg.SearchDataFile(terraformWorkingDir)
	if err != nil {
		log.Fatalf("error getting Terraform working directory: %s", err)
	}

	return workingDir
}

func getTerraformInstallDir() string {
	installDir, err := xdg.CacheFile(terraformInstallDir)
	if err != nil {
		log.Fatalf("error creating Terraform install directory: %s", err)
	}

	return installDir
}

func getTerraformExecutable() string {
	installer := install.NewInstaller()
	installDirectory := getTerraformInstallDir()
	versionConstraint := version.MustConstraints(version.NewConstraint(terraformVersion))

	executable, err := installer.Ensure(context.Background(), []src.Source{
		&fs.Version{
			Product:     product.Terraform,
			ExtraPaths:  []string{installDirectory},
			Constraints: versionConstraint,
		},
		&releases.LatestVersion{
			Product:     product.Terraform,
			InstallDir:  installDirectory,
			Constraints: versionConstraint,
		},
	})

	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	return executable
}
