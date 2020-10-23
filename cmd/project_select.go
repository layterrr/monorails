package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func selectProject(name string) error {
	projectsConfig, err := readProjectsConfig()
	if err != nil {
		return err
	}
	projectsConfig.Selected = name
	if err := updateProjectsConfig(projectsConfig); err != nil {
		return err
	}
	return nil
}

var selectProjectCmd = &cobra.Command{
	Use:     "select",
	Aliases: []string{"workon", "switch", "set"},
	Short:   "Select a project",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Must specify a project name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := selectProject(name); err != nil {
			panic(err)
		}
		fmt.Printf("Selected project %s", name)
	},
}

func init() {
	projectCmd.AddCommand(selectProjectCmd)
}
