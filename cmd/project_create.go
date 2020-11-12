package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

func createProject(name string) error {
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
	projectsConfig.Projects[name] = path.Join(pwd, name)
	if err := updateProjectsConfig(projectsConfig); err != nil {
		return err
	}
	return nil
}

var createProjectCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "init"},
	Short:   "Create a new project",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Must specify a project name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := createProject(name); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		if err := selectProject(name); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created project %s\n", name)
	},
}

func init() {
	projectCmd.AddCommand(createProjectCmd)
}
