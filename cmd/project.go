package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// projectTemplateRepo is repo used to generate new monorails projects
var projectTemplateRepo = "https://github.com/layterrr/gcp-project-template"

// projectsConfig represents the monorails projects stored in the monorails.yml file
type projectsConfig struct {
	configPath string
	Projects   map[string]string `yaml:"projects"`
	Selected   string            `yaml:"selected"`
}

// selectedProject returns the currently selected project
func (c *projectsConfig) selectedProject() string {
	return c.Projects[c.Selected]
}

func (c *projectsConfig) selectProject(project string) error {
	c.Selected = project
	return c.update()
}

func (c *projectsConfig) addProject(projectName, projectPath string) error {
	if _, ok := c.Projects[projectName]; ok {
		return errors.New("Project already exists")
	}
	c.Projects[projectName] = projectPath
	return c.update()
}

func (c *projectsConfig) removeProject(project string) error {
	if _, ok := c.Projects[project]; !ok {
		return errors.New("Project doesn't exist")
	}
	delete(c.Projects, project)
	return c.update()
}

// update writes the config object to the monorails.yml file
func (c *projectsConfig) update() error {
	if c.configPath == "" {
		return errors.New("No config file set")
	}

	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	f, err := os.Create(c.configPath)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(b); err != nil {
		return err
	}
	return nil
}

// newProjectsConfig reads project config file and creates a new
// projectsConfig struct
func newProjectsConfig() (*projectsConfig, error) {
	c := &projectsConfig{
		Projects: map[string]string{},
	}

	// Lookup custom config path and default to ~/.monorails.yml
	projectsConfig, ok := os.LookupEnv("MONORAILS_CONFIG")
	if !ok {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		projectsConfig = path.Join(home, ".monorails.yml")
	}
	c.configPath = projectsConfig

	// Read config file and instantiate a new projectsConfig
	in, err := ioutil.ReadFile(c.configPath)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(in, c); err != nil {
		return nil, err
	}

	return c, nil
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
	Aliases: []string{"projects"},
	Short:   "Group command for managing projects",
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.PersistentFlags().StringVarP(
		&projectTemplateRepo,
		"template",
		"t",
		projectTemplateRepo,
		"Project template repo url",
	)
}
