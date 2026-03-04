package internal

import (
	"context"
	"log"
	"net"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-version"
	install "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hc-install/src"
	"github.com/hashicorp/terraform-exec/tfexec"
)

type Terraform struct {
	instance *tfexec.Terraform
}

type IngressRule struct {
	Cidr string `json:"cidr_ipv4"`
	Port uint   `json:"port"`
}

type TerraformVariables struct {
	Profile        string `json:"aws_profile"`
	Region         string `json:"aws_region"`
	KeyDirectory   string `json:"key_directory"`
	DeploymentId   string `json:"deployment_id"`
	DeploymentType string `json:"deployment_type"`
	AdminCidr      string `json:"admin_cidr,omitempty"`

	IngressRules []IngressRule `json:"ingress_rules"`
}

type TerraformOutput struct {
	DeploymentId  string
	InstanceIp    string
	SshKeyFile    string
	CloudFrontUrl string
}

const (
	defaultWorkspace = "default"
	terraformVersion = "~> 1.10"
)

func (*Terraform) New() *Terraform {
	StartSpinner("Initializing RedSky")

	tf, err := tfexec.NewTerraform(getTerraformWorkingDir(), getTerraformExecutable())
	if err != nil {
		log.Fatalf("error creating Terraform instance: %s", err)
	}

	if err := tf.Init(context.Background(), tfexec.Backend(false)); err != nil {
		log.Fatalf("error initializing Terraform: %s", err)
	}

	StopSpinner()

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

	StopSpinner()

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

func (tf *Terraform) ApplyDeployment(profile string, region string, deploymentType string, adminCidr net.IPNet, ingressRules []IngressRule) {
	StartSpinner("Executing deployment")

	workspaceName := tf.CreateWorkspace()

	vars := TerraformVariables{
		Profile:        profile,
		Region:         region,
		KeyDirectory:   getRedSkyDir(),
		DeploymentId:   workspaceName,
		DeploymentType: deploymentType,
	}

	for _, ingressRule := range ingressRules {
		vars.IngressRules = append(vars.IngressRules, ingressRule)
	}

	if adminCidr.IP != nil {
		vars.AdminCidr = adminCidr.String()
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

	StopSpinner()
}

func (tf *Terraform) UpdateDeployment(workspaceName string, ingressRules []IngressRule) {
	StartSpinner("Updating deployment")

	if err := tf.instance.WorkspaceSelect(context.Background(), workspaceName); err != nil {
		log.Fatalf("error selecting Terraform workspace: %s", err)
	}

	varFilePath := getVarFilePath(workspaceName)

	var vars TerraformVariables
	if err := readVarsFile(varFilePath, &vars); err != nil {
		log.Fatalf("error reading tfvars file: %s", err)
	}

	for _, ingressRule := range ingressRules {
		vars.IngressRules = append(vars.IngressRules, ingressRule)
	}

	if err := writeVarFile(varFilePath, vars); err != nil {
		StopSpinnerError("Deployment failed")
		log.Fatalf("error updating tfvars file: %s", err)
	}

	if err := tf.instance.Apply(context.Background(), tfexec.VarFile(varFilePath)); err != nil {
		StopSpinnerError("Deployment failed")
		log.Fatalf("error deploying: %s", err)
	}

	StopSpinner()
}

func (tf *Terraform) DestroyDeployment(workspaceName string) {
	StartSpinner("Destroying deployment")

	if err := tf.instance.WorkspaceSelect(context.Background(), workspaceName); err != nil {
		log.Fatalf("error selecting Terraform workspace: %s", err)
	}

	if err := tf.instance.Destroy(context.Background(), tfexec.VarFile(getVarFilePath(workspaceName))); err != nil {
		log.Fatalf("error destroying Terraform deployment: %s", err)
	}

	tf.DeleteWorkspace(workspaceName)

	StopSpinner()
}

func (tf *Terraform) GetDeploymentDetails() *TerraformOutput {
	StartSpinner("Gathering deployment details")

	outputs, err := tf.instance.Output(context.Background())
	if err != nil {
		log.Fatalf("error retrieving Terraform output: %s", err)
	}

	StopSpinner()

	return &TerraformOutput{
		DeploymentId:  strings.Trim(string(outputs["deployment_id"].Value), "\""),
		InstanceIp:    strings.Trim(string(outputs["instance_ip"].Value), "\""),
		SshKeyFile:    strings.Trim(string(outputs["ssh_key_file"].Value), "\""),
		CloudFrontUrl: strings.Trim(string(outputs["cloudfront_url"].Value), "\""),
	}
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
