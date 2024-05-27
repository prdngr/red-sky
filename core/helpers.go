package core

import (
	"log"
	"os"
	"os/user"
	"path"
)

const (
	nodDirectory         = ".nod"
	directoryPermissions = 0755
)

func GetNodDirectory() string {
	user, err := user.Current()
	if err != nil {
		log.Fatalf("error getting current user: %s", err)
	}

	return path.Join(user.HomeDir, nodDirectory)
}

func DirectoryExists(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	} else {
		return true, err
	}
}

func CreateDirectoryIfNotExists(path string) error {
	if exists, err := DirectoryExists(path); err != nil {
		return err
	} else if !exists {
		if err := os.MkdirAll(path, directoryPermissions); err != nil {
			return err
		}
	}

	return nil
}
