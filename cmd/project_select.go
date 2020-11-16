package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func selectProject(name string) error {
	config, err := newProjectsConfig()
	if err != nil {
		return err
	}

	if err := config.selectProject(name); err != nil {
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
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Selected project %s\n", name)
	},
}

func init() {
	projectCmd.AddCommand(selectProjectCmd)
}
