package cmd

import (
	"errors"
	"fmt"
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

func startServices(services []string) error {
	args := []string{
		"-f", "./docker-compose.yml",
	}

	for _, service := range services {
		args = append(args, "-f")
		args = append(args, fmt.Sprintf("%s/docker-compose.yml", service))
	}
	args = append(args, "up")
	args = append(args, "--build")
	cmd := exec.Command("docker-compose", args...)
	cmd.Dir = synthwaveDir
	fmt.Println(synthwaveDir, cmd.Dir)
	_, err := runCommand("Running docker-compose up", cmd)
	if err != nil {
		return err
	}
	return nil
}

var startServicesCmd = &cobra.Command{
	Use:   "start [services]",
	Short: "Start services",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 && !allServices {
			return errors.New("You must define at least one service to start or use the --all flag")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		serviceNames := args
		if allServices {
			serviceNames, err = listServices()
			if err != nil {
				log.Fatal("Error listing services:", err)
			}
		}

		if err := startServices(serviceNames); err != nil {
			log.Fatal("Error starting services:", err)
		}
	},
}

func init() {
	serviceCmd.AddCommand(startServicesCmd)
}
