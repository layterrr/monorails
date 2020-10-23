package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func applyCommand() error {
	checkEnvironment("apply")
	currentProjectID := getProjectID()

	projectName, err := getProjectName(currentProjectID)
	if err != nil {
		return err
	}
	formattedProjectName := formatProjectName(projectName)

	vars, err := projectVars(projectName, formattedProjectName, currentProjectID)
	if err != nil {
		return err
	}

	printCmd("Generate terraform apply command. Run this in ./infra",
		terraformCommand([]string{"apply"}, vars))

	return nil
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Generate a terraform apply command",
	Long:  `Copy this command and run in ./infra`,
	Run: func(cmd *cobra.Command, args []string) {
		err := applyCommand()
		if err != nil {
			fmt.Printf("Error generating terraform command: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	environmentCmd.AddCommand(applyCmd)
}
