package cmd

import (
	"errors"
	"log"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

func checkService(service string) error {
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
	cmd.Dir = path.Join(synthwaveDir, service)
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
				panic(err)
			}
		}

		for _, service := range services {
			err := checkService(service)
			if err != nil {
				log.Fatal("Failed to check service:", err.Error())
			}
		}
	},
}

func init() {
	serviceCmd.AddCommand(checkServiceCmd)
}
