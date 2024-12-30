package core

import (
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/adrg/xdg"
	"github.com/prdngr/nessus-on-demand/static"
)

const (
	NodDir          = "nod/"
	initializedFile = NodDir + ".initialized"
	filePermissions = 0660
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
	if _, err := xdg.SearchDataFile(initializedFile); err == nil {
		return
	}

	if err := fs.WalkDir(static.Embeds, ".", func(path string, entry fs.DirEntry, err error) error {
		if entry.IsDir() {
			return nil
		}

		if data, err := static.Embeds.ReadFile(path); err != nil {
			return err
		} else {
			if outputPath, err := xdg.DataFile(NodDir + path); err != nil {
				return err
			} else {
				if err := os.WriteFile(outputPath, data, filePermissions); err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		log.Fatalf("error initializing NoD directory: %s", err)
	}
}

func PrintBanner() {
	fmt.Fprintln(os.Stderr, nodBanner)
}
