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
	Short: "Show health check for the gorevive command.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := exec.LookPath("tmux")
		if err != nil {
			fmt.Println("tmux is installed? No")
		} else {
			fmt.Println("tmux is installed? Yes")
		}

		editor := os.Getenv("EDITOR")
		if editor != "" {
			fmt.Printf("$EDITOR is set: yes\n")
		} else {
			fmt.Println("$EDITOR is set: no")
		}
		shell := os.Getenv("SHELL")
		if shell != "" {
			fmt.Printf("$SHELL is set: yes\n")
		} else {
			fmt.Println("$SHELL is set: no")

		}
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
}
