package cmd

import (
	"errors"
	"log"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

func testService(service string) error {
	serviceConfig, err := readServiceConfig(service)
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
	cmd.Dir = path.Join(synthwaveDir, service)
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
			services, err = listServices()
			if err != nil {
				panic(err)
			}
		}

		for _, service := range services {
			err := testService(service)
			if err != nil {
				log.Fatal("Failed to test service:", err.Error())
			}
		}
	},
}

func init() {
	serviceCmd.AddCommand(testServiceCmd)
}
