package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func shipService(service *serviceConfig) (err error) {
	if err := testService(service.Path); err != nil {
		return err
	}
	if err := checkService(service.Path); err != nil {
		return err
	}
	if err := genProto(service.Path); err != nil {
		return err
	}
	if err := buildService(service); err != nil {
		return err
	}
	if err := deployService(service); err != nil {
		return err
	}
	return nil
}

// shipCmd represents the ship command
var shipCmd = &cobra.Command{
	Use:   "ship",
	Short: "Tests, build, and deploys a service",
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
			fmt.Println(strings.Join(serviceNames, ", "))
			serviceConfig, err := readServiceConfig(serviceName)
			if err != nil {
				fmt.Printf("Error reading service config: %v\n", err)
				os.Exit(1)
			}

			err = shipService(serviceConfig)
			if err != nil {
				fmt.Printf("Error shipping service: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	serviceCmd.AddCommand(shipCmd)
}
