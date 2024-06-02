package helpers

import (
	"fmt"
	"os/exec"
	"strings"
)

// Command encapsulates the execution of a command.
type Command struct {
	Parts []string
}

// Add appends parts to the command.
func (m *Command) Add(part ...string) {
	m.Parts = append(m.Parts, part...)
}

// Execute runs the command with the given base and arguments.
func (m *Command) Execute(base string, args []string) (string, error) {
	args = append(args, m.Parts...)
	cmd := exec.Command(base, args...)

	var outBuffer, errBuffer strings.Builder
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute %s %s: %s", base, strings.Join(args, " "), errBuffer.String())
	}
	return strings.TrimSpace(outBuffer.String()), nil
}

// Clear resets the command parts.
func (m *Command) Clear() {
	m.Parts = nil
}
