package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var defaultRegion = "europe-west1"

var environmentCmd = &cobra.Command{
	Use:     "environment",
	Aliases: []string{"environments", "env"},
	Short:   "Commands to create and manage gcloud environments",
}

func init() {
	rootCmd.AddCommand(environmentCmd)
}

func getProjectID() string {
	if isTest {
		return environments["test"]
	}
	if isStaging {
		return environments["staging"]
	}
	if isProd {
		return environments["production"]
	}
	return environments["test"]
}

func getProjectNumber(projectID string) (string, error) {
	cmd := exec.Command("gcloud", "projects", "describe", projectID)
	out, err := runCommand("Get project number", cmd)
	if err != nil {
		return "", err
	}

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		kv := strings.Split(strings.TrimSpace(line), ":")
		if kv[0] == "projectNumber" {
			number := strings.TrimSpace(kv[1])
			return strings.ReplaceAll(number, "'", ""), nil
		}
	}

	return "", errors.New("Couldn't find project")
}

func getProjectName(projectID string) (string, error) {
	cmd := exec.Command("gcloud", "projects", "describe", projectID)
	out, err := runCommand("Get project name", cmd)
	if err != nil {
		return "", err
	}

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		kv := strings.Split(strings.TrimSpace(line), ":")
		if kv[0] == "name" {
			name := strings.ReplaceAll(kv[1], "'", "")
			return strings.TrimSpace(name), nil
		}
	}

	return "", errors.New("Couldn't find project")
}

func generateProjectID(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, randomString(6))
}

func formatProjectName(projectName string) string {
	return strings.ReplaceAll(strings.ToLower(projectName), " ", "-")
}

func projectVars(projectName, formattedProjectName, projectID string) (map[string]string, error) {
	vars := map[string]string{}
	vars["project_id"] = projectID
	vars["project_name"] = projectName
	return vars, nil
}

func appendUniqueSecrets(secrets map[string]string) (map[string]string, error) {
	return secrets, nil
}

func terraformCommand(args []string, vars map[string]string) (*exec.Cmd, error) {
	cmd := exec.Command("terraform", args...)
	projectsConfig, err := newProjectsConfig()
	if err != nil {
		return nil, err
	}
	projectDir := projectsConfig.selectedProject()
	infraDir := path.Join(projectDir, "infra")
	cmd.Dir = infraDir

	for k, v := range vars {
		tfVar := fmt.Sprintf("TF_VAR_%s=\"%s\"", k, v)
		cmd.Env = append(cmd.Env, tfVar)
	}

	return cmd, nil
}

func createTerraformWorkspace(name string) error {
	initTerraformCmd, err := terraformCommand([]string{"init"}, nil)
	if err != nil {
		return err
	}
	if _, err := runCommand("Initialise terraform workspace", initTerraformCmd); err != nil {
		return err
	}

	createTerraformWorkspaceCmd, err := terraformCommand([]string{"workspace", "new", name}, nil)
	if err != nil {
		return err
	}
	if _, err := runCommand("Create new terraform workspace", createTerraformWorkspaceCmd); err != nil {
		return err
	}

	return nil
}

func checkEnvironment(commandName string) {
	if !isTest && !force {
		fmt.Printf("Use the `--force` flag to run %s against a non test environment\n", commandName)
	}
}
