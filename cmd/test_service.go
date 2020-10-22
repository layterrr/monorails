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

func testService(service string) error {
	serviceConfig, err := readServiceConfig("projectDir", service)
	if err != nil {
		return err
	}

	if serviceConfig.TestCommand == "" {
		return nil
	}

	commandParts := strings.Split(serviceConfig.TestCommand, " ")
	commandName := commandParts[0]
	commandArgs := commandParts[1:]
	cmd := exec.Command(commandName, commandArgs...)
	cmd.Dir = path.Join("projectDir", service)
	_, err = runCommand("Testing service", cmd)

	if err != nil {
		return err
	}

	return nil
}

var testServiceCmd = &cobra.Command{
	Use:   "test [service]",
	Short: "Run tests for the service",
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
			services, err = listServices("projectDir")
			if err != nil {
				fmt.Printf("Failed to list services: %v\n", err)
				os.Exit(1)
			}
		}

		for _, service := range services {
			err := testService(service)
			if err != nil {
				fmt.Printf("Failed to test service: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	serviceCmd.AddCommand(testServiceCmd)
}
