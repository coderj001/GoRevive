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
		showNumbers, _ := cmd.Flags().GetInt("number")

		if err != nil {
			fmt.Println(err)
		}
		for idx, cfg := range cfgs {
			if idx >= showNumbers {
				return
			}
			fmt.Println(cfg)
		}

	},
}

func init() {
	listCmd.Flags().IntP("number", "n", 10, "Number projects the list")
	rootCmd.AddCommand(listCmd)
}
