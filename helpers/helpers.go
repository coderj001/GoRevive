package helpers

import (
	"fmt"
	"os"
	"path/filepath"
)

// Get Current Configuration Dir
func getCurrentConfigDir() string {
	user, err := os.Current()
	if err != nil {
		panic(err)
	}
	return filepath.Join(user.HomeDir, ".config", "gorevive")
}

// check or create for given path
func CheckOrCreateDir(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(path) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return false, fmt.Errorf("failed to create directory: %w", err)
		}
	} else {
		return true, nil
	}
}

// GetFile returns the path to an existing file with given filename
func GetFile(filename string) (string, error) {}

// CreateFile returns the path to an file after creating the file
func CreateFile(filename string) (string, error) {}

// Get all config files list
func GetConfigFiles() ([]string, error) {

}
