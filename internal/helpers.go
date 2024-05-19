package internal

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
)

func GetNodDirectory() string {
	user, err := user.Current()
	if err != nil {
		log.Fatalf("Error getting current user: %s", err)
	}

	return path.Join(user.HomeDir, NOD_DIRECTORY)
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
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("error creating directory: %w", err)
		}
	}

	return nil
}
