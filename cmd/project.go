package cmd

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type projectsConfig struct {
	Projects map[string]string `yaml:"projects"`
	Selected string            `yaml:"selected"`
}

func projectsConfigPath() (string, error) {
	projectConfig, ok := os.LookupEnv("MONORAILS_CONFIG")
	if ok {
		return projectConfig, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(home, ".monorails.yml"), nil
}

func readProjectsConfig() (*projectsConfig, error) {
	config := &projectsConfig{
		Projects: map[string]string{},
	}
	configPath, err := projectsConfigPath()
	if err != nil {
		return nil, err
	}
	in, err := ioutil.ReadFile(configPath)
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
	configPath, err := projectsConfigPath()
	if err != nil {
		return err
	}
	f, err := os.Create(configPath)
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
