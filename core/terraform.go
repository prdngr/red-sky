package core

import (
	"context"
	"log"

	"github.com/adrg/xdg"
	"github.com/hashicorp/go-version"
	install "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hc-install/src"
	"github.com/hashicorp/terraform-exec/tfexec"
)

const (
	versionConstraint   = "~> 1.8"
	terraformWorkingDir = "terraform/"
	terraformInstallDir = "bin/"
)

func GetTerraformInstance() *tfexec.Terraform {
	tf, err := tfexec.NewTerraform(getTerraformWorkingDir(), getTerraformExecutable())
	if err != nil {
		log.Fatalf("error creating Terraform instance: %s", err)
	}

	if err := tf.Init(context.Background(), tfexec.Backend(false)); err != nil {
		log.Fatalf("error initializing Terraform: %s", err)
	}

	return tf
}

func getTerraformWorkingDir() string {
	workingDir, err := xdg.SearchDataFile(NodDir + terraformWorkingDir)
	if err != nil {
		log.Fatalf("error getting Terraform working directory: %s", err)
	}

	return workingDir
}

func getTerraformInstallDir() string {
	installDir, err := xdg.CacheFile(NodDir + terraformInstallDir + "x")
	if err != nil {
		log.Fatalf("error creating Terraform install directory: %s", err)
	}

	return installDir
}

func getTerraformExecutable() string {
	installer := install.NewInstaller()
	installDirectory := getTerraformInstallDir()
	versionConstraint := version.MustConstraints(version.NewConstraint(versionConstraint))

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
