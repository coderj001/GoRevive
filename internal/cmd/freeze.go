package cmd

import (
	"github.com/coderj001/GoRevive/internal/tmux"
	"github.com/spf13/cobra"
)

// freezeCmd represents the freeze command
var freezeCmd = &cobra.Command{
	Use:   "freeze [project]",
	Short: "Freeze an existing project",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		tmux.FreezeSession()
	},
}

func init() {
	rootCmd.AddCommand(freezeCmd)
}
