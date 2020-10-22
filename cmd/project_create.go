package cmd

import (
	"github.com/spf13/cobra"
)

func createProject() ([]string, error) {
	return nil, nil
}

var createProjectCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "init"},
	Short:   "Create a new project",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	projectCmd.AddCommand(createProjectCmd)
}
