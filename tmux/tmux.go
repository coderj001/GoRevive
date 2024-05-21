package tmux

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/coderj001/GoRevive/helpers"
	"gopkg.in/yaml.v2"
)

type Project struct {
	Name    string   `yaml:"project_name,omitempty"`
	Root    string   `yaml:"project_root,omitempty"`
	OnStart []string `yaml:"on_project_start,omitempty"`
	OnEnd   []string `yaml:"on_project_end,omitempty"`
}

// Session handle a tmux session, each session contain many Window
type Session struct {
	Attach  bool     `yaml:"-"`
	Name    string   `yaml:"name,omitempty"`
	Windows []Window `yaml:"windows,omitempty"`
}

// Window handle tmux window, each window can have multiple pane.
type Window struct {
	Name     string        `yaml:"name"`
	Root     string        `yaml:"root,omitempty"`
	Layout   string        `yaml:"layout"`
	Panes    []interface{} `yaml:"panes"`
	RealPane []Pane        `yaml:"-"`
}

// Pane handle each pane (single command line) in tmux
type Pane struct {
	Commands   []string `yaml:"commands"`
	Focus      bool     `yaml:"focus,omitempty"`
	Root       string   `yaml:"root,omitempty"`
	identifier string   `yaml:"-"`
}

// Command is a helper for executable command inside tmux pane
type Command struct {
	Parts []string
}

var (
	// IgnoredCmd is list of commands ignored by save session
	IgnoredCmd []string
	// DefaultCmd is a command used when the command is ignored
	DefaultCmd string
	// Copy of env whitout tmux related env
	tmuxENV []string
)

// func BuildSession() {}

// CreateSession new tmux sessions.
func CreateSession(sessionName string) error {
	if IsInsideTmux() {
		return fmt.Errorf("already inside tmux session")
	}
	_, err := helpers.RunCommand("tmux", "new-session", "-s", sessionName)
	return err
}

// IsInsideTmux Check if we are inside tmux or not
func IsInsideTmux() bool {
	// Simply, if the TMUX is set in env, We are in it :)
	return os.Getenv("TMUX") != ""
}

// ListSessions lists all the tmux sessions.
func ListSessions() ([]string, error) {
	output, err := helpers.RunCommand("tmux", "list-sessions", "-F", "#S")
	if err != nil {
		return nil, err
	}
	sessionNames := strings.Split(strings.TrimSpace(output), "\n")
	return sessionNames, nil
}

func BuildSession() error {
	filePath := "/home/mrzero/.config/gorevive/config.yaml"
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var p Project
	err = yaml.Unmarshal(content, &p)
	if err != nil {
		return fmt.Errorf("failed to unmarshal YAML content: %w", err)
	}

	fmt.Printf("----------- x -----------\n")
	fmt.Printf("--- Project:\n%v\n\n", p)
	return nil
}
