package internal

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/adrg/xdg"
	"github.com/prdngr/red-sky/static"
)

const (
	redSkyDir           = "red-sky/"
	redSkyVersionFile   = redSkyDir + ".version"
	terraformWorkingDir = redSkyDir + "terraform/"
	filePermissions     = 0700
)

func InitRedSkyDir() {
	toolVersion := GetVersion().BuildVersion
	installedVersion := readVersionFile()

	if toolVersion == installedVersion && toolVersion != "dev" {
		return
	}

	if err := fs.WalkDir(static.Embeds, ".", func(path string, entry fs.DirEntry, err error) error {
		if entry.IsDir() {
			return nil
		}

		data, err := static.Embeds.ReadFile(path)
		if err != nil {
			return err
		}

		if fileNeedsUpdate(path, data) {
			diskPath, err := xdg.DataFile(redSkyDir + path)
			if err != nil {
				return err
			}

			if err := os.WriteFile(diskPath, data, filePermissions); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.Fatalf("error initializing RedSky directory: %s", err)
	}

	writeVersionFile(toolVersion)
}

func getRedSkyDir() string {
	directory, err := xdg.SearchDataFile(redSkyDir)
	if err != nil {
		log.Fatalf("error getting RedSky directory: %s", err)
	}

	return directory
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

func readVersionFile() string {
	if versionFile, err := xdg.SearchDataFile(redSkyVersionFile); err == nil {
		if data, err := os.ReadFile(versionFile); err == nil {
			return strings.TrimSpace(string(data))
		}
	}

	return ""
}

func writeVersionFile(version string) {
	versionFile, err := xdg.DataFile(redSkyVersionFile)
	if err != nil {
		log.Fatalf("error writing version file: %s", err)
	}

	os.WriteFile(versionFile, []byte(version), filePermissions)
}

func fileNeedsUpdate(path string, data []byte) bool {
	diskPath, err := xdg.DataFile(redSkyDir + path)
	if err != nil {
		return true
	}

	existingData, err := os.ReadFile(diskPath)
	if err != nil {
		return true
	}

	return sha256.Sum256(existingData) != sha256.Sum256(data)
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
