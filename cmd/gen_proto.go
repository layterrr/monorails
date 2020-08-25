package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func genProto(service string) error {
	serviceConfig, err := readServiceConfig(service)
	if err != nil {
		return err
	}

	if serviceConfig.GenProtoCommand == "" {
		return nil
	}

	cmd := createCommand(serviceConfig.GenProtoCommand, synthwaveDir)
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
				panic(err)
			}
		}

		for _, serviceName := range serviceNames {
			err = genProto(serviceName)
			if err != nil {
				log.Fatal("Failed to generate proto:", err.Error())
			}
		}
	},
}

func init() {
	serviceCmd.AddCommand(genProtoCmd)
}
