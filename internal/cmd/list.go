package cmd

import (
	"fmt"

	"github.com/coderj001/GoRevive/internal/helpers"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long:  `List all projects configs`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgs, err := helpers.GetConfigFiles()

		if err != nil {
			fmt.Println(err)
		}
		for _, cfg := range cfgs {
			fmt.Println(cfg)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
