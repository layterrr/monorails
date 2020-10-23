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
			panic(err)
		}
		if err := selectProject(name); err != nil {
			panic(err)
		}
		fmt.Printf("Created project %s", name)
	},
}

func init() {
	projectCmd.AddCommand(createProjectCmd)
}
