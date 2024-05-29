package cmd

import (
	"fmt"

	"github.com/coderj001/GoRevive/internal/helpers"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [project]",
	Short: "Delete an existing project",
	Long:  `Delete a project configuration file specified by the project name.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		err := helpers.DeleteFile(filename)
		if err != nil {
			fmt.Printf("Error deleting file %s: %v\n", filename, err)
		} else {
			fmt.Printf("Successfully deleted file %s\n", filename)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
