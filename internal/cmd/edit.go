package cmd

import (
	"fmt"

	"github.com/coderj001/GoRevive/internal/helpers"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [project]",
	Short: "Edit an existing project",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		err := helpers.EditFile(filename)
		if err != nil {
			fmt.Printf("Error editing file %s: %v\n", filename, err)
		} else {
			fmt.Printf("Successfully editing file %s\n", filename)
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
