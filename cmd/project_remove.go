package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func removeProject(name string) error {
	projectsConfig, err := readProjectsConfig()
	if err != nil {
		return err
	}
	delete(projectsConfig.Projects, name)
	if projectsConfig.Selected == name {
		projectsConfig.Selected = ""
	}
	if err := updateProjectsConfig(projectsConfig); err != nil {
		return err
	}
	return nil
}

// projectRemoveCmd represents the projectRemove command
var removeProjectCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"delete", "rm"},
	Short:   "Removes project from monorails config. It will not remove the repository.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Must specify at least one project name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		names := args

		for _, name := range names {
			if err := removeProject(name); err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
		}
		fmt.Printf("Removed projects: %s.\n", strings.Join(names, ", "))
		fmt.Println("Note: this has not deleted any files")
	},
}

func init() {
	projectCmd.AddCommand(removeProjectCmd)
}
