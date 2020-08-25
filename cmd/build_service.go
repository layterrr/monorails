package cmd

import (
	"errors"
	"fmt"
	"log"
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
				panic(err)
			}
		}

		for _, serviceName := range serviceNames {
			serviceConfig, err := readServiceConfig(serviceName)
			if err != nil {
				panic(err)
			}

			err = buildService(serviceConfig)
			if err != nil {
				log.Fatal("Failed to build service:", err.Error())
			}
		}
	},
}

func init() {
	serviceCmd.AddCommand(buildServiceCmd)
}
