package cmd

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var version string
var allServices bool

var serviceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"services", "svc"},
	Short:   "Commands for creating and managing services",
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.PersistentFlags().StringVarP(&version, "version", "v", "latest", "version")
	serviceCmd.PersistentFlags().BoolVarP(&allServices, "all", "a", false, "all services")
}

type serviceConfig struct {
	Name            string `yaml:"name"`
	Path            string
	TestCommand     string   `yaml:"test-command"`
	CheckCommand    string   `yaml:"check-command"`
	GenProtoCommand string   `yaml:"gen-proto-command"`
	IncludePaths    []string `yaml:"include"`
}

func readServiceConfig(service string) (*serviceConfig, error) {
	projectsConfig, err := newProjectsConfig()
	if err != nil {
		return nil, err
	}
	projectDir := projectsConfig.selectedProject()

	config := &serviceConfig{}
	in, err := ioutil.ReadFile(path.Join(projectDir, service, "service.yml"))
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(in, config); err != nil {
		return nil, err
	}

	config.Path = service

	if len(config.Name) == 0 {
		return nil, fmt.Errorf("Service name is required")
	}

	return config, nil
}
