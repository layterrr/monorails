package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listProjectsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := readProjectsConfig()
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		if len(config.Projects) == 0 {
			fmt.Println("No projects found")
		}
		for project := range config.Projects {
			line := ""
			if project == config.Selected {
				line = "* "
			}
			line = line + project
			fmt.Println(line)
		}
	},
}

func init() {
	projectCmd.AddCommand(listProjectsCmd)
}
