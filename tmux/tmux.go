package tmux

import (
	"strings"

	"github.com/coderj001/GoRevive/helpers"
)

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
	_, err := helpers.RunCommand("tmux", "new-session", "-s", sessionName)
	return err
}

// func AttachSession(sessionName string) error { }

// ListSessions lists all the tmux sessions.
func ListSessions() ([]string, error) {
	output, err := helpers.RunCommand("tmux", "list-sessions", "-F", "#S")
	if err != nil {
		return nil, err
	}
	sessionNames := strings.Split(strings.TrimSpace(output), "\n")
	return sessionNames, nil
}
