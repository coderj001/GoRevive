// helpers
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

func NewFile(project string) error {
	filePath, err := createFile(project)
	if err != nil {
		return err
	}
	err = insertDefaultConfig(project, filePath)
	if err != nil {
		return err
	}
	return nil
}

func insertDefaultConfig(project, filePath string) error {
	content := fmt.Sprintf(`project_name: %s
# project_root: ~/src/project_path
# on_project_start:
#   - sudo systemctl start postgresql
# pre_window:
#   - workon dummy
# windows:
#   - editor: vim
#   - shells:
#       layout: main-vertical
#       panes:
#         - #
#         - grunt serve`, project)

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write default config to %s: %w", filePath, err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	cmd := exec.Command(editor, filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("unable to edit file %s: %w", filePath, err)
	}
	return nil
}

// createFile creates a new file and returns its path
func createFile(filename string) (string, error) {
	filePath := filepath.Join(getCurrentConfigDir(), fmt.Sprintf("%s.yaml", filename))

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return "", fmt.Errorf("file %s already exists", filePath)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()
	return filePath, nil
}

// DeleteFile deletes the specified configuration file.
func DeleteFile(filename string) error {
	path := filepath.Join(getCurrentConfigDir(),
		fmt.Sprintf("%s.yaml", filename))

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", path)
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("unable to delete file %s: %w", path, err)
	}
	return nil
}

// EditFile edit the specified configuration file.
func EditFile(filename string) error {
	path := filepath.Join(getCurrentConfigDir(),
		fmt.Sprintf("%s.yaml", filename))

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", path)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	cmd := exec.Command(editor, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("unable to edit file %s: %w", path, err)
	}
	return nil
}

// GetConfigFiles fetch all config files list
func GetConfigFiles() ([]string, error) {
	// TODO: sort out base of filters
	configDirPath := getCurrentConfigDir()
	ok, err := checkOrCreatePath(configDirPath, true)
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

	cmd.Stdin = os.Stdin

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()

	output := outBuf.String() + errBuf.String()

	return output, err
}

// checkOrCreatePath checks if a path exists. If it doesn't, it creates the path as a directory or file based on the provided isDir flag.
func checkOrCreatePath(path string, isDir bool) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if isDir {
			if err := os.MkdirAll(path, 0755); err != nil {
				return false, fmt.Errorf("failed to create directory: %w", err)
			}
		} else {
			dir := filepath.Dir(path)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return false, fmt.Errorf("failed to create directory: %w", err)
				}
			}

			file, err := os.Create(path)
			if err != nil {
				return false, fmt.Errorf("failed to create file: %w", err)
			}
			defer file.Close()
		}
		return true, nil
	}
	return true, nil
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
