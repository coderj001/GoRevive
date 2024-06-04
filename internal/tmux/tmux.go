package tmux

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/coderj001/GoRevive/internal/helpers"
	"gopkg.in/yaml.v2"
)

type Project struct {
	Name    string   `yaml:"project_name,omitempty"`
	Root    string   `yaml:"project_root,omitempty"`
	OnStart []string `yaml:"on_project_start,omitempty"`
	OnEnd   []string `yaml:"on_project_end,omitempty"`
	Windows []Window `yaml:"windows,omitempty"`
}

// Session handles a tmux session, each session contains many Window.
type Session struct {
	Attach  bool     `yaml:"-"`
	Name    string   `yaml:"name,omitempty"`
	Windows []Window `yaml:"windows,omitempty"`
}

// Window handles tmux windows, each window can have multiple panes.
type Window struct {
	Name   string `yaml:"name"`
	Index  int    `yaml:"-"`
	Height int    `yaml:"height"`
	Width  int    `yaml:"width"`
	Panes  []Pane `yaml:"panes"`
	Focus  bool   `yaml:"focus,omitempty"`
}

// Pane handles each pane (single command line) in tmux.
type Pane struct {
	Commands []string `yaml:"commands"`
	Index    int      `yaml:"-"`
	Focus    bool     `yaml:"focus,omitempty"`
	Height   int      `yaml:"height"`
	Width    int      `yaml:"width"`
}

var cmd = &helpers.Command{}

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
		"#{window_index} #{window_name} #{window_width} #{window_height} #{window_active}",
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
		active := parts[4] == "1"
		windows = append(windows,
			Window{Index: index,
				Name:   parts[1],
				Height: height,
				Width:  width,
				Focus:  active,
			})
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
		parts := strings.Split(line, " ")
		index, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("invalid pane index: %v", err)
		}
		height, err := strconv.Atoi(parts[2])
		if err != nil {
			return fmt.Errorf("invalid window height: %v", err)
		}
		width, err := strconv.Atoi(parts[3])
		if err != nil {
			return fmt.Errorf("invalid window width: %v", err)
		}
		active := parts[1] == "1"
		panes = append(panes, Pane{
			Index:    index,
			Focus:    active,
			Height:   height,
			Width:    width,
			Commands: []string{parts[4]},
		})
	}
	window.Panes = panes
	return nil
}

func (s *Session) setProject(pr *Project) error {
	out, err := cmd.Execute("tmux", []string{
		"display-message",
		"-p",
		"-F",
		"#{pane_current_path}",
		"-t",
		s.Name,
	})
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
	_, err := cmd.Execute("tmux", []string{"new-session", "-s", sessionName})
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

	err = helpers.CreateFile(project.Name, data)
	if err != nil {
		return fmt.Errorf("error creating YAML file: %v", err)
	}

	return nil
}

// ListSessions lists all the tmux sessions.
func ListSessions() ([]string, error) {
	output, err := cmd.Execute("tmux", []string{"list-sessions", "-F", "#S"})
	if err != nil {
		return nil, err
	}
	sessionNames := strings.Split(strings.TrimSpace(output), "\n")
	return sessionNames, nil
}

func BuildSession(sessionName string) error {
	if IsInsideTmux() {
		return fmt.Errorf("already inside tmux session")
	}

	data, err := helpers.LoadData(sessionName)
	if err != nil {
		return fmt.Errorf("Unable to fetch file: %w", err)
	}

	var project Project
	err = yaml.Unmarshal(*data, &project)
	if err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	sessionName = project.Name

	cmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create tmux session: %w", err)
	}

	fmt.Printf("Created new tmux session: %s\n", sessionName)

	for _, window := range project.Windows {
		windowCmd := exec.Command("tmux", "new-window", "-t", sessionName, "-n", window.Name)
		err = windowCmd.Run()
		if err != nil {
			return fmt.Errorf("failed to create tmux window: %w", err)
		}

		for _, pane := range window.Panes {
			for _, command := range pane.Commands {
				paneCmd := exec.Command("tmux", "split-window", "-t", fmt.Sprintf("%s:%s", sessionName, window.Name), "-h", command)
				err = paneCmd.Run()
				if err != nil {
					return fmt.Errorf("failed to create tmux pane: %w", err)
				}
			}
		}
	}

	for _, command := range project.OnStart {
		onStartCmd := exec.Command("tmux", "send-keys", "-t", sessionName, command, "C-m")
		err = onStartCmd.Run()
		if err != nil {
			return fmt.Errorf("failed to execute on_start command: %w", err)
		}
	}

	fmt.Printf("Configured tmux session: %s\n", sessionName)
	return nil
}
