package cmd

import (
	"fmt"
	"log"
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
				panic(err)
			}

			err = shipService(serviceConfig)
			if err != nil {
				log.Fatal("Failed to ship service:", err.Error())
			}
		}
	},
}

func init() {
	serviceCmd.AddCommand(shipCmd)
}
