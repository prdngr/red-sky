package core

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/adrg/xdg"
	"github.com/fatih/color"
	"github.com/prdngr/nessus-on-demand/internal"
	"github.com/prdngr/nessus-on-demand/static"
)

const (
	NodDir          = "nod/"
	initializedFile = NodDir + ".version"
	filePermissions = 0700
	nodBanner       = `
    ███╗   ██╗ ██████╗ ██████╗
    ████╗  ██║██╔═══██╗██╔══██╗
    ██╔██╗ ██║██║   ██║██║  ██║
    ██║╚██╗██║██║   ██║██║  ██║
    ██║ ╚████║╚██████╔╝██████╔╝
    ╚═╝  ╚═══╝ ╚═════╝ ╚═════╝
	`
)

func GetNodDir() string {
	nodDir, err := xdg.SearchDataFile(NodDir)
	if err != nil {
		log.Fatalf("error getting NoD directory: %s", err)
	}

	return nodDir
}

func InitNodDir() {
	NodVersion := internal.GetVersion().BuildVersion
	IacVersion := readVersionFile()

	if NodVersion == IacVersion && NodVersion != "dev" {
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
			diskPath, err := xdg.DataFile(NodDir + path)
			if err != nil {
				return err
			}

			if err := os.WriteFile(diskPath, data, filePermissions); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.Fatalf("error updating NoD directory: %s", err)
	}

	writeVersionFile(NodVersion)
}

func PrintBanner() {
	fmt.Fprintln(os.Stderr, nodBanner)
}

func PrintHeader(header string) {
	color.Yellow("\n" + header)
	color.Yellow(strings.Repeat("-", len(header)) + "\n\n")
}

func readVersionFile() string {
	if versionFile, err := xdg.SearchDataFile(initializedFile); err == nil {
		if data, err := os.ReadFile(versionFile); err == nil {
			return strings.TrimSpace(string(data))
		}
	}

	return ""
}

func writeVersionFile(version string) {
	versionFile, err := xdg.DataFile(initializedFile)
	if err != nil {
		log.Fatalf("error writing version file: %s", err)
	}

	os.WriteFile(versionFile, []byte(version), filePermissions)
}

func fileNeedsUpdate(path string, data []byte) bool {
	diskPath, err := xdg.DataFile(NodDir + path)
	if err != nil {
		return true
	}

	existingData, err := os.ReadFile(diskPath)
	if err != nil {
		return true
	}

	return sha256.Sum256(existingData) != sha256.Sum256(data)
}
