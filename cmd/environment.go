package cmd

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/spf13/cobra"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var defaultRegion = "europe-west1"
var sharedSecrets = map[string]bool{}

var environmentCmd = &cobra.Command{
	Use:     "environment",
	Aliases: []string{"environments"},
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

func projectVars(projectName, formattedProjectName, projectID string, secrets map[string]string) (map[string]string, error) {
	vars := secrets
	vars["project_id"] = projectID
	vars["project_name"] = projectName
	return vars, nil
}

func readSecrets(projectID string, secretKeys []string) (map[string]string, error) {
	secrets := map[string]string{}

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to create secretmanager client: %v", err)
	}

	for _, secret := range secretKeys {
		name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secret)
		req := &secretmanagerpb.AccessSecretVersionRequest{
			Name: name,
		}

		result, err := client.AccessSecretVersion(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to access secret version: %v", err)
		}

		secretValue := string(result.Payload.Data)
		secrets[secret] = strings.TrimSpace(secretValue)
	}

	return secrets, nil
}

func readSharedSecrets(projectID string) (map[string]string, error) {
	secretKeys := []string{}
	for key, shared := range sharedSecrets {
		if shared {
			secretKeys = append(secretKeys, key)
		}
	}
	return readSecrets(projectID, secretKeys)
}

func readAllSecrets(projectID string) (map[string]string, error) {
	secretKeys := []string{}
	for key := range sharedSecrets {
		secretKeys = append(secretKeys, key)
	}
	return readSecrets(projectID, secretKeys)
}

func appendUniqueSecrets(secrets map[string]string) (map[string]string, error) {
	return secrets, nil
}

func terraformCommand(args []string, vars map[string]string) *exec.Cmd {
	cmd := exec.Command("terraform", args...)
	cmd.Dir = infraDir

	for k, v := range vars {
		tfVar := fmt.Sprintf("TF_VAR_%s=\"%s\"", k, v)
		cmd.Env = append(cmd.Env, tfVar)
	}

	return cmd
}

func createTerraformWorkspace(name string) error {
	if _, err := runCommand("Create new terraform workspace", terraformCommand([]string{"workspace", "new", name}, nil)); err != nil {
		return err
	}
	if _, err := runCommand("Initialise terraform workspace", terraformCommand([]string{"init"}, nil)); err != nil {
		return err
	}
	return nil
}

func checkEnvironment(commandName string) {
	if !isTest && !force {
		fmt.Printf("Use the `--force` flag to run %s against a non test environment\n", commandName)
	}
}
