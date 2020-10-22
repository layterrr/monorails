package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func listProjects() ([]string, error) {
	return nil, nil
}

var listProjectsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Run: func(cmd *cobra.Command, args []string) {
		projects, err := listProjects()
		if err != nil {
			panic(err)
		}
		for _, project := range projects {
			fmt.Println(project)
		}
	},
}

func init() {
	projectCmd.AddCommand(listProjectsCmd)
}
