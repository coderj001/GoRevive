// helpers
package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/coderj001/GoRevive/internal/config"
)

func NewFile(project string) error {
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
	err := CreateFile(project, []byte(content))
	if err != nil {
		return err
	}
	return nil
}

func LoadData(fileName string) (*[]byte, error) {
	filepath := filepath.Join(
		config.ConfigDir,
		fmt.Sprintf("%s.yaml", fileName),
	)
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", fileName, err)
	}

	return &data, nil
}

func extractFileName(fileName string) *string {
	parts := strings.Split(fileName, ".")
	if len(parts) >= 2 {
		return &parts[0]
	}
	return nil
}

func checkIfPathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return true
	}
	return false
}
