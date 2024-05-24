package tmux

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/coderj001/GoRevive/helpers"
	"gopkg.in/yaml.v2"
)

type Project struct {
	Name    string   `yaml:"project_name,omitempty"`
	Root    string   `yaml:"project_root,omitempty"`
	OnStart []string `yaml:"on_project_start,omitempty"`
	OnEnd   []string `yaml:"on_project_end,omitempty"`
	Windows []Window `yaml:"windows,omitempty"`
}

// Session handle a tmux session, each session contain many Window
type Session struct {
	Attach  bool     `yaml:"-"`
	Name    string   `yaml:"name,omitempty"`
	Windows []Window `yaml:"windows,omitempty"`
}

// Window handle tmux window, each window can have multiple pane.
type Window struct {
	Name   string `yaml:"name"`
	Index  int    `yaml:"-"`
	Height int    `yaml:"height"`
	Width  int    `yaml:"width"`
	Panes  []Pane `yaml:"panes"`
	Focus  bool   `yaml:"focus,omitempty"`
}

// Pane handle each pane (single command line) in tmux
type Pane struct {
	Commands []string `yaml:"commands"`
	Index    int      `yaml:"-"`
	Focus    bool     `yaml:"focus,omitempty"`
	Height   int      `yaml:"height"`
	Width    int      `yaml:"width"`
}

var (
	// IgnoredCmd is list of commands ignored by save session
	IgnoredCmd []string
	// DefaultCmd is a command used when the command is ignored
	DefaultCmd string
	// Copy of env whitout tmux related env
	tmuxENV []string

	cmd = &Command{}
)

type Command struct {
	Parts []string
}

func (m *Command) Add(part ...string) {
	if m.Parts == nil {
		m.Parts = make([]string, 0)
	}
	m.Parts = append(m.Parts, part...)
}

func (m *Command) Execute(base string, args []string) (string, error) {
	args = append(args, m.Parts...)
	cmd := exec.Command(base, args...)
	cmd.Env = tmuxENV

	var outBuffer, errBuffer strings.Builder
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	err := cmd.Run()

	if err != nil {
		err = fmt.Errorf("failed to execute %s %s : %s \n",
			base,
			strings.Join(args, " "),
			errBuffer.String())
	}
	return strings.TrimSpace(outBuffer.String()), err
}

func (m *Command) Clear() {
	m.Parts = nil
}

func (s *Session) getSessionNames() []string {
	out, err := cmd.Execute("tmux",
		[]string{
			"list-sessions",
			"-F",
			"#S",
		})
	if err != nil {
		fmt.Println(err)
	}
	sessionNames := strings.Split(strings.TrimSpace(out), "\n")
	return sessionNames
}

func (s *Session) setSelectedSession() error {
	sessions := s.getSessionNames()
	if len(sessions) > 1 {
		fmt.Printf("Select from range [%d-%d]:\n", 0, len(sessions)-1)
		for i, sessionName := range sessions {
			fmt.Printf("[%d] %s\n", i, sessionName)
		}

		var sessionIndex int
		fmt.Print(" ===> ")
		_, err := fmt.Scanf("%d", &sessionIndex)
		if err != nil {
			return fmt.Errorf("failed to read input: %v", err)
		}

		if sessionIndex < 0 || sessionIndex >= len(sessions) {
			return fmt.Errorf("invalid selection: %d", sessionIndex)
		}

		fmt.Println("Selected session:", sessions[sessionIndex])
		s.Name = sessions[sessionIndex]
	} else if len(sessions) == 1 {
		s.Name = sessions[0]
		fmt.Println("Only one session available, selected:", s.Name)
	} else {
		return fmt.Errorf("no sessions available")
	}
	return nil
}

func (s *Session) setCurrentSession() error {
	if IsInsideTmux() {
		out, err := cmd.Execute("tmux", []string{"display-message", "-p", "#S"})
		if err != nil {
			return fmt.Errorf("error getting sessions: %v", err)
		}
		s.Name = out
	} else {
		err := s.setSelectedSession()
		if err != nil {
			return fmt.Errorf("error getting sessions: %v", err)
		}

	}
	return nil
}

func (s *Session) setWindows() error {
	out, err := cmd.Execute("tmux", []string{
		"list-windows",
		"-t",
		s.Name,
		"-F",
		"#{window_index} #{window_name} #{window_width} #{window_height}",
	})
	if err != nil {
		return fmt.Errorf("error getting windows: %v", err)
	}

	lines := strings.Split(out, "\n")
	windows := make([]Window, 0, len(lines))
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}
		index, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("invalid window index: %v", err)
		}
		height, err := strconv.Atoi(parts[2])
		if err != nil {
			return fmt.Errorf("invalid window height: %v", err)
		}
		width, err := strconv.Atoi(parts[3])
		if err != nil {
			return fmt.Errorf("invalid window width: %v", err)
		}
		windows = append(windows,
			Window{Index: index,
				Name:   parts[1],
				Height: height,
				Width:  width})
	}
	s.Windows = windows
	return nil
}

func (s *Session) setPanes(window *Window) error {
	out, err := cmd.Execute("tmux",
		[]string{
			"list-panes",
			"-t",
			fmt.Sprintf("%s:%d", s.Name, window.Index),
			"-F",
			"#{pane_index} #{pane_active} #{pane_height} #{pane_width} #{pane_current_command}",
		})
	if err != nil {
		return fmt.Errorf("error getting panes: %v", err)
	}

	lines := strings.Split(out, "\n")
	panes := make([]Pane, 0, len(lines))
	for _, line := range lines {
		l := strings.Split(line, " ")
		index, err := strconv.Atoi(l[0])
		if err != nil {
			return fmt.Errorf("invalid pane index: %v", err)
		}
		height, err := strconv.Atoi(l[2])
		if err != nil {
			return fmt.Errorf("invalid window height: %v", err)
		}
		width, err := strconv.Atoi(l[3])
		if err != nil {
			return fmt.Errorf("invalid window width: %v", err)
		}
		active := l[1] == "1" // Convert "1" to true and "0" to false
		panes = append(panes, Pane{
			Index:    index,
			Focus:    active,
			Height:   height,
			Width:    width,
			Commands: []string{l[4]},
		})
	}
	window.Panes = panes
	return nil
}

func (s *Session) setProject(pr *Project) error {
	out, err := cmd.Execute("pwd", nil)
	if err != nil {
		return err
	}
	pr.Name = s.Name
	pr.Root = out
	pr.Windows = s.Windows
	return nil
}

// IsInsideTmux Check if we are inside tmux or not
func IsInsideTmux() bool {
	return os.Getenv("TMUX") != ""
}

// CreateSession new tmux sessions.
func CreateSession(sessionName string) error {
	if IsInsideTmux() {
		return fmt.Errorf("already inside tmux session")
	}
	_, err := helpers.RunCommand("tmux", "new-session", "-s", sessionName)
	return err
}

func FreezeSession() error {
	s := &Session{}
	project := &Project{}
	err := s.setCurrentSession()
	if err != nil {
		return err
	}
	err = s.setWindows()
	if err != nil {
		return err
	}
	for i := range s.Windows {
		err := s.setPanes(&s.Windows[i])
		if err != nil {
			return err
		}
	}
	s.setProject(project)

	data, err := yaml.Marshal(project)
	if err != nil {
		return fmt.Errorf("error marshaling project to YAML: %v", err)
	}

	file, err := os.Create("project.yaml")
	if err != nil {
		return fmt.Errorf("error creating YAML file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing to YAML file: %v", err)
	}

	return nil
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
