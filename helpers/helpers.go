package helpers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CheckOrCreateDir checks if a directory exists at the given path, creates it if it doesn't exist, and returns true if it exists or is created successfully.
func CheckOrCreateDir(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return false, fmt.Errorf("failed to create directory: %w", err)
		}
		return true, nil
	}
	return true, nil
}

// GetFile returns the path to an existing file with given filename
// func GetFile(fileName string) (string, error){ }

// CreateFile returns the path to an file after creating the file
func CreateFile(filename string) error {
	file, err := os.Create(fmt.Sprintf("%w.yaml", filename))
	if err != nil {
		return fmt.Errorf("failed to create file %w", err)
	}
	defer file.Close()
	return nil
}

// DeleteFile
func DeleteFile(filename string) error {
	path := filepath.Join(getCurrentConfigDir(), filename)
	err := os.Remove(path)
	return fmt.Errorf("Unable to create, %w", err)
}

// GetConfigFiles fetch all config files list
func GetConfigFiles() ([]string, error) {
	// TODO: sort out base of filters
	configDirPath := getCurrentConfigDir()
	ok, err := CheckOrCreateDir(configDirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get config dir %w", err)
	}
	if ok {
		files, err := ioutil.ReadDir(configDirPath)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch files %w", err)
		}
		var configFiles []string
		for _, file := range files {
			fileName := extractFileName(file.Name())
			if fileName != nil {
				configFiles = append(configFiles, *fileName)
			}
		}
		return configFiles, nil
	}
	return nil, fmt.Errorf("No Files Found.")
}

// RunCommand executes a shell command and returns its output or an error.
func RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	return out.String(), err
}

func getCurrentConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, ".config", "gorevive")
}

func extractFileName(fileName string) *string {
	parts := strings.Split(fileName, ".")
	if len(parts) >= 2 {
		return &parts[0]
	}
	return nil
}
