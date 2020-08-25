package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/cobra"
)

func listServices() ([]string, error) {
	services := []string{}
	servicesDir := path.Join(projectDir, "services")
	files, err := ioutil.ReadDir(servicesDir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			servicePath := path.Join("services", file.Name())
			serviceYamlPath := path.Join(servicePath, "service.yml")
			_, err := os.Stat(serviceYamlPath)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return nil, err
			}
			services = append(services, servicePath)
		}
	}
	return services, nil
}

var listServicesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all services",
	Run: func(cmd *cobra.Command, args []string) {
		services, err := listServices()
		if err != nil {
			panic(err)
		}
		for _, service := range services {
			fmt.Println(service)
		}
	},
}

func init() {
	serviceCmd.AddCommand(listServicesCmd)
}
