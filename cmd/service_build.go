package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func buildService(service *serviceConfig) (err error) {
	_, err = runCommand("Build service", exec.Command(
		"gcloud",
		"builds",
		"submit",
		fmt.Sprintf("--config=%s/build.yml", service.Path),
		fmt.Sprintf("--substitutions=_VERSION=%s", version),
		fmt.Sprintf("--project=%s", containerProjectID),
	))
	return err
}

var buildServiceCmd = &cobra.Command{
	Use:   "build [service]",
	Short: "Build a container for the service",
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
				fmt.Printf("Error listing services: %v", err)
				os.Exit(1)
			}
		}

		for _, serviceName := range serviceNames {
			serviceConfig, err := readServiceConfig(serviceName)
			if err != nil {
				fmt.Printf("Error reading service config: %v", err)
				os.Exit(1)
			}

			err = buildService(serviceConfig)
			if err != nil {
				fmt.Printf("Error building service: %v", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	serviceCmd.AddCommand(buildServiceCmd)
}
