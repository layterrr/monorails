package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func addProject(name string) error {
	config, err := newProjectsConfig()
	if err != nil {
		return err
	}

	if err := config.addProject(name); err != nil {
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
