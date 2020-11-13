package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func deployService(service *serviceConfig) error {
	currentProjectID := getProjectID()
	_, err := runCommand("Deploy service", exec.Command(
		"gcloud",
		"builds",
		"submit",
		"--no-source",
		fmt.Sprintf("--config=%s/deploy.yml", service.Path),
		fmt.Sprintf("--project=%s", currentProjectID),
		fmt.Sprintf("--substitutions=_VERSION=%s,_CONTAINER_PROJECT=%s", version, containerProjectID),
	))

	if err != nil {
		return err
	}

	return nil
}

var deployServiceCmd = &cobra.Command{
	Use:   "deploy [service]",
	Short: "Deploy a version of the service",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 && !allServices {
			return errors.New("service required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		serviceNames := args
		if allServices {
			serviceNames, err = listServices()
			if err != nil {
				panic(err)
			}
		}

		for _, serviceName := range serviceNames {
			serviceConfig, err := readServiceConfig(serviceName)
			if err != nil {
				fmt.Printf("Error reading service config: %v", err)
				os.Exit(1)
			}

			err = deployService(serviceConfig)
			if err != nil {
				fmt.Printf("Error deploying service: %v", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	serviceCmd.AddCommand(deployServiceCmd)
}
