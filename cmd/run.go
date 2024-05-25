package cmd

import (
	"fmt"

	"github.com/coderj001/GoRevive/tmux"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [project]",
	Short: "run an existing project",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please provide a project name.")
			return
		}
		project := args[0]
		err := tmux.BuildSession(project)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
