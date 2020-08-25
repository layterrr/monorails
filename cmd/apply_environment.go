package cmd

import (
	"log"

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

	secrets, err := readAllSecrets(currentProjectID)
	if err != nil {
		return err
	}

	vars, err := projectVars(projectName, formattedProjectName, currentProjectID, secrets)
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
			log.Fatal("Failed to generate terraform apply command:", err.Error())
		}
	},
}

func init() {
	environmentCmd.AddCommand(applyCmd)
}
