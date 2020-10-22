package cmd

import (
	"github.com/spf13/cobra"
)

func selectProject() ([]string, error) {
	return nil, nil
}

var selectProjectCmd = &cobra.Command{
	Use:     "select",
	Aliases: []string{"workon", "switch"},
	Short:   "Select a project",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	projectCmd.AddCommand(selectProjectCmd)
}
