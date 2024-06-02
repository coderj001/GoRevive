// helpers files operations
package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/coderj001/GoRevive/internal/config"
)

// CreateFile creates a new file and returns its path
func CreateFile(filename string, content []byte) error {
	filePath := filepath.Join(
		config.ConfigDir,
		fmt.Sprintf("%s.yaml", filename),
	)

	if checkIfPathExists(filePath) {
		return fmt.Errorf("file %s already exists", filePath)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}

	return nil
}

// DeleteFile deletes the specified configuration file.
func DeleteFile(filename string) error {
	filePath := filepath.Join(config.ConfigDir,
		fmt.Sprintf("%s.yaml", filename))

	if checkIfPathExists(filePath) {
		return fmt.Errorf("file %s does not exist", filePath)
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("unable to delete file %s: %w", filePath, err)
	}
	return nil
}

// EditFile edit the specified configuration file.
func EditFile(filename string) error {
	filePath := filepath.Join(
		config.ConfigDir,
		fmt.Sprintf("%s.yaml", filename),
	)
	fmt.Println(filePath)

	if checkIfPathExists(filePath) {
		return fmt.Errorf("file %s does not exist", filePath)
	}

	editorStr := os.Getenv("EDITOR")
	if editorStr == "" {
		return fmt.Errorf("editor is not set", filePath)
	}

	editor, err := exec.LookPath(editorStr)
	if err != nil {
		return err
	}
	if err := syscall.Exec(
		editor,
		[]string{editorStr, filePath},
		os.Environ()); err != nil {
		return err
	}
	return nil
}

// GetConfigFiles fetch all config files list
func GetConfigFiles() ([]string, error) {
	files, err := ioutil.ReadDir(config.ConfigDir)
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
