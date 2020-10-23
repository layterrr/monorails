package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func listProjects() ([]string, error) {
	config, err := readProjectsConfig()
	if err != nil {
		return nil, err
	}
	projects := []string{}
	for project := range config.Projects {
		projects = append(projects, project)
	}
	return projects, nil
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
