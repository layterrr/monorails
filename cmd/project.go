package cmd

import (
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type projectsConfig struct {
	Projects map[string]string `yaml:"projects"`
	Selected string            `yaml:"selected"`
}

func projectsConfigPath() string {
	projectConfig, ok := os.LookupEnv("MONORAILS_CONFIG")
	if ok {
		return projectConfig
	}
	return "~/.monorails.yaml"
}

func readProjectsConfig() (*projectsConfig, error) {
	config := &projectsConfig{}
	in, err := ioutil.ReadFile(projectsConfigPath())
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(in, config); err != nil {
		return nil, err
	}
	return config, nil
}

func updateProjectsConfig(config *projectsConfig) error {
	b, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	f, err := os.Create(projectsConfigPath())
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(b); err != nil {
		return err
	}
	return nil
}

var projectCmd = &cobra.Command{
	Use:     "project",
	Short:   "Group command for managing projects",
	Aliases: []string{"projects"},
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
