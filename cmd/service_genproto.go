package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func genProto(service string) error {
	projectDir, err := getProjectDirectory()
	if err != nil {
		return err
	}
	serviceConfig, err := readServiceConfig(service)
	if err != nil {
		return err
	}

	if serviceConfig.GenProtoCommand == "" {
		return nil
	}

	cmd := createCommand(serviceConfig.GenProtoCommand, projectDir)
	_, err = runCommand(fmt.Sprintf("Generating protos for %s", serviceConfig.Name), cmd)
	if err != nil {
		return err
	}

	return nil
}

var genProtoCmd = &cobra.Command{
	Use:   "gen-proto [service]",
	Short: "Generate protos for a service",
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
			err = genProto(serviceName)
			if err != nil {
				fmt.Printf("Error generating protos: %v", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	serviceCmd.AddCommand(genProtoCmd)
}
