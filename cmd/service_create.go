package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var serviceTemplate string

type serviceTemplateVars struct {
	ServiceName string
	ProjectDir  string
	ServicesDir string
}

func createService(name string) error {
	if serviceTemplate == "" {
		return errors.New("No template provided")
	}

	projectsConfig, err := readProjectsConfig()
	if err != nil {
		return err
	}

	projectDir := projectsConfig.Projects[projectsConfig.Selected]
	serviceDir := path.Join(projectDir, "services", name)

	if err := os.Mkdir(serviceDir, 0755); err != nil {
		if os.IsExist(err) {
			return errors.New("Service already exists")
		}
		return err
	}

	vars := &serviceTemplateVars{
		ServiceName: name,
		ProjectDir:  projectDir,
		ServicesDir: path.Join(projectDir, "services"),
	}

	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
		"Title":   strings.Title,
	}

	isRoot := true
	err = filepath.Walk(serviceTemplate, func(templatePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Hack to prevent creating a directory of the template root
		if isRoot {
			isRoot = false
			return nil
		}

		if info.IsDir() {
			if err := os.Mkdir(path.Join(serviceDir, info.Name()), 0755); err != nil {
				return err
			}
			return nil
		}

		t, err := template.New(info.Name()).Funcs(funcMap).ParseFiles(templatePath)
		if err != nil {
			return err
		}

		strippedFilePath := strings.Replace(templatePath, serviceTemplate, "", 1)
		outputPath := path.Join(serviceDir, strippedFilePath)
		outputPath = strings.Replace(outputPath, ".tmpl", "", 1)
		f, err := os.Create(outputPath)
		if err != nil {
			return err
		}

		if err := t.Execute(f, vars); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

var createServiceCmd = &cobra.Command{
	Use:     "create [service]",
	Aliases: []string{"new", "gen"},
	Short:   "Create a new service from a template",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("No service name provided")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		if err := createService(name); err != nil {
			fmt.Printf("Error creating service: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	serviceCmd.AddCommand(createServiceCmd)
	createServiceCmd.Flags().StringVarP(&serviceTemplate, "template", "t", "", "Template to generate service from")
}
