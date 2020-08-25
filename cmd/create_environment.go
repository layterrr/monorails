package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var duplicateProjectID string

func createEnvironment(projectName string) (string, error) {
	checkEnvironment("create")

	formattedProjectName := formatProjectName(projectName)

	if err := createTerraformWorkspace(formattedProjectName); err != nil {
		return "", err
	}

	projectID := generateProjectID(formattedProjectName)
	secrets, err := readSharedSecrets(duplicateProjectID)
	if err != nil {
		return "", err
	}

	secrets, err = appendUniqueSecrets(secrets)
	if err != nil {
		return "", err
	}

	vars, err := projectVars(projectName, formattedProjectName, projectID, secrets)
	if err != nil {
		return "", err
	}

	printCmd("Generate terraform plan command. Run this in ./infra",
		terraformCommand([]string{"apply"}, vars))
	return projectID, nil
}

var createEnvironmentCmd = &cobra.Command{
	Use:     "create [project-name]",
	Aliases: []string{"new"},
	Short:   "Create a new environment on gcloud",
	Long:    `Copy this command and run in ./infra`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("project_name required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		_, err := createEnvironment(projectName)
		if err != nil {
			fmt.Printf("Error creating environment: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	environmentCmd.AddCommand(createEnvironmentCmd)
	createEnvironmentCmd.PersistentFlags().StringVarP(&duplicateProjectID,
		"duplicate", "d", stagingProject, "Project to duplicate")
}
