package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
)

func startServices(services []string) error {
	projectsConfig, err := newProjectsConfig()
	if err != nil {
		return err
	}
	projectDir := projectsConfig.selectedProject()

	args := []string{
		"-f", path.Join(projectDir, "docker-compose.yml"),
	}

	for _, service := range services {
		args = append(args, "-f")
		args = append(args, path.Join(projectDir, service, "/docker-compose.yml"))
	}
	args = append(args, "up")
	args = append(args, "--build")
	cmd := exec.Command("docker-compose", args...)
	cmd.Dir = projectDir
	if _, err := runCommand("Running docker-compose up", cmd); err != nil {
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
				fmt.Printf("Error listing services: %v\n", err)
				os.Exit(1)
			}
		}

		if err := startServices(serviceNames); err != nil {
			fmt.Printf("Error starting services: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	serviceCmd.AddCommand(startServicesCmd)
}
