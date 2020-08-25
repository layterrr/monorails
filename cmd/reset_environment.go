package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func resetEnvironment() (string, error) {
	checkEnvironment("reset")
	currentProjectID := getProjectID()

	projectName, err := getProjectName(currentProjectID)
	if err != nil {
		return "", err
	}
	formattedProjectName := formatProjectName(projectName)

	newProjectID := generateProjectID(formattedProjectName)
	secrets, err := readSharedSecrets(currentProjectID)
	if err != nil {
		return "", err
	}

	secrets, err = appendUniqueSecrets(secrets)
	if err != nil {
		return "", err
	}

	vars, err := projectVars(projectName, formattedProjectName, newProjectID, secrets)
	if err != nil {
		return "", err
	}

	printCmd("Generate terraform plan command. Run this in ./infra",
		terraformCommand([]string{"apply"}, vars))
	return newProjectID, nil
}

var resetEnvironmentCmd = &cobra.Command{
	Use:   "reset",
	Short: "Resets an existing environment on gcloud",
	Long: `Copy this command and run in ./infra
	This will replace the project in the current terraform
	workspace.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := resetEnvironment()
		if err != nil {
			fmt.Printf("Error resetting environment: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	environmentCmd.AddCommand(resetEnvironmentCmd)
}
