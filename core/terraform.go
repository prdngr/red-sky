package core

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

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
	terraformVersion    = "~> 1.10"
	terraformWorkingDir = NodDir + "terraform/"
)

type Terraform struct {
	instance *tfexec.Terraform
}

type TerraformOutput struct {
	DeploymentId string
	InstanceIp   string
	SshKeyFile   string
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
	workspaceName := uuid.NewString()

	if err := tf.instance.WorkspaceNew(context.Background(), workspaceName); err != nil {
		log.Fatalf("error creating Terraform workspace: %s", err)
	}

	return workspaceName
}

func (tf *Terraform) DeleteWorkspace(workspace string) {
	if err := tf.instance.WorkspaceSelect(context.Background(), defaultWorkspace); err != nil {
		log.Fatalf("error selecting Terraform workspace: %s", err)
	}

	if err := tf.instance.WorkspaceDelete(context.Background(), workspace); err != nil {
		log.Fatalf("error deleting Terraform workspace: %s", err)
	}
}

func (tf *Terraform) ApplyDeployment(profile string, region string, deploymentType string, allowedIp net.IP) {
	StartSpinner("Deploying Nessus")

	workspaceName := tf.CreateWorkspace()

	vars := map[string]string{
		"aws_profile":     profile,
		"aws_region":      region,
		"key_directory":   GetNodDir(),
		"deployment_id":   workspaceName,
		"deployment_type": deploymentType,
	}

	if allowedIp.To4() != nil && !allowedIp.IsLoopback() {
		vars["allowed_ip"] = allowedIp.String()
	}

	varFilePath := getVarFilePath(workspaceName)

	if err := writeVarFile(varFilePath, vars); err != nil {
		StopSpinnerError("Deployment failed")
		tf.DeleteWorkspace(workspaceName)
		log.Fatalf("error creating tfvars file: %s", err)
	}

	if err := tf.instance.Apply(context.Background(), tfexec.VarFile(varFilePath)); err != nil {
		StopSpinnerError("Deployment failed")
		tf.DeleteWorkspace(workspaceName)
		log.Fatalf("error deploying: %s", err)
	}

	StopSpinner("Nessus deployed")
}

func (tf *Terraform) DestroyDeployment(profile string, workspaceName string) {
	StartSpinner("Destroying deployment")

	if err := tf.instance.WorkspaceSelect(context.Background(), workspaceName); err != nil {
		log.Fatalf("error selecting Terraform workspace: %s", err)
	}

	if err := tf.instance.Destroy(context.Background(), tfexec.VarFile(getVarFilePath(workspaceName))); err != nil {
		log.Fatalf("error destroying Terraform deployment: %s", err)
	}

	tf.DeleteWorkspace(workspaceName)

	StopSpinner("Deployment destroyed")
}

func (tf *Terraform) GetDeploymentDetails() *TerraformOutput {
	StartSpinner("Gathering deployment details")

	outputs, err := tf.instance.Output(context.Background())
	if err != nil {
		log.Fatalf("error retrieving Terraform output: %s", err)
	}

	StopSpinner("Deployment details gathered")

	return &TerraformOutput{
		DeploymentId: strings.Trim(string(outputs["deployment_id"].Value), "\""),
		InstanceIp:   strings.Trim(string(outputs["instance_ip"].Value), "\""),
		SshKeyFile:   strings.Trim(string(outputs["ssh_key_file"].Value), "\""),
	}
}

func writeVarFile(varFilePath string, vars map[string]string) error {
	varFile, err := os.Create(varFilePath)
	if err != nil {
		return err
	}
	defer varFile.Close()

	for key, value := range vars {
		if _, err = fmt.Fprintf(varFile, "%s = \"%s\"\n", key, value); err != nil {
			return err
		}
	}

	return nil
}

func getVarFilePath(workspaceName string) string {
	varFile, err := xdg.DataFile(terraformWorkingDir + workspaceName + ".tfvars")
	if err != nil {
		log.Fatalf("error searching for tfvars file: %s", err)
	}

	return varFile
}

func getTerraformWorkingDir() string {
	workingDir, err := xdg.SearchDataFile(terraformWorkingDir)
	if err != nil {
		log.Fatalf("error getting Terraform working directory: %s", err)
	}

	return workingDir
}

func getTerraformInstallDir() string {
	terraformInstallDir := xdg.BinHome

	if err := os.MkdirAll(terraformInstallDir, filePermissions); err != nil {
		log.Fatalf("error creating Terraform install directory: %s", err)
	}

	return terraformInstallDir
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
