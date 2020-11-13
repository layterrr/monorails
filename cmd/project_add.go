package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

func addProject(name string) error {
	projectsConfig, err := readProjectsConfig()
	if err != nil {
		return err
	}

	if _, ok := projectsConfig.Projects[name]; ok {
		return errors.New("Project already exists")
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	projectDir := path.Join(pwd, name)
	projectsConfig.Projects[name] = projectDir
	if err := updateProjectsConfig(projectsConfig); err != nil {
		return err
	}

	return nil
}

var addProjectCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"register"},
	Short:   "Adds an existing repo as a monorails project",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Must specify a project to add")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := addProject(name); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Added project: %s.\n", name)
	},
}

func init() {
	projectCmd.AddCommand(addProjectCmd)
}
