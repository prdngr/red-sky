package core

import (
	"context"
	"log"
	"path"

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
	workingDirectory    = "terraform"
	executableDirectory = "bin"
)

func InstallTerraform() {
	installer := install.NewInstaller()
	versionConstraint := version.MustConstraints(version.NewConstraint(versionConstraint))
	installDirectory := path.Join(GetNodDirectory(), executableDirectory)

	CreateDirectoryIfNotExists(installDirectory)

	executablePath, err := installer.Ensure(context.Background(), []src.Source{
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

	Config.Terraform.ExecutablePath = executablePath
}

func InitializeTerraform() *tfexec.Terraform {
	workingDirectory := path.Join(GetNodDirectory(), workingDirectory)
	executablePath := Config.Terraform.ExecutablePath

	tf, err := tfexec.NewTerraform(workingDirectory, executablePath)
	if err != nil {
		log.Fatalf("error creating Terraform instance: %s", err)
	}

	err = tf.Init(context.Background())
	if err != nil {
		log.Fatalf("error initializing Terraform: %s", err)
	}

	return tf
}
