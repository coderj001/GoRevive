package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Show health check for the gorevive command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		healthCheck()
	},
}

func healthCheck() {
	_, err := exec.LookPath("tmux")
	if err != nil {
		fmt.Println("tmux is installed? No")
	} else {
		fmt.Println("tmux is installed? Yes")
	}

	editor := os.Getenv("EDITOR")
	if editor != "" {
		fmt.Printf("$EDITOR is set: Yes\n")
	} else {
		fmt.Printf("$EDITOR is set: No\n")
	}
	shell := os.Getenv("SHELL")
	if shell != "" {
		fmt.Printf("$SHELL is set: Yes\n")
	} else {
		fmt.Printf("$SHELL is set: No\n")

	}
}

func init() {
	rootCmd.AddCommand(healthCmd)
}
