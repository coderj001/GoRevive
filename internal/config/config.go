package config

import (
	"log"
	"os"
	"path/filepath"
)

// ConfigDir get config directory
var ConfigDir string

func init() {
	userDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error: could not determine user home: %v", err)
	}

	ConfigDir = filepath.Join(userDir, ".config", "gorevive")

	if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
		if err := os.MkdirAll(ConfigDir, 0755); err != nil {
			log.Fatalf("error: failed to create directory: %w", err)
		}
	}

}
