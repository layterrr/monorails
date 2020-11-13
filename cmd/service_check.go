package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

func checkService(service string) error {
	projectDir, err := getProjectDirectory()
	if err != nil {
		return err
	}
	serviceConfig, err := readServiceConfig(service)
	if err != nil {
		return err
	}

	if serviceConfig.TestCommand == "" {
		return nil
	}

	commandParts := strings.Split(serviceConfig.CheckCommand, " ")
	commandName := commandParts[0]
	commandArgs := commandParts[1:]
	cmd := exec.Command(commandName, commandArgs...)
	cmd.Dir = path.Join(projectDir, service)
	_, err = runCommand("Checking service", cmd)

	if err != nil {
		return err
	}

	return nil
}

var checkServiceCmd = &cobra.Command{
	Use:   "check [service]",
	Short: "Run checks for the service",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("service required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		services := args
		if allServices {
			services, err = listServices()
			if err != nil {
				fmt.Printf("Error listing services: %v", err)
				os.Exit(1)
			}
		}

		for _, service := range services {
			err := checkService(service)
			if err != nil {
				fmt.Printf("Error checking service: %v", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	serviceCmd.AddCommand(checkServiceCmd)
}
