package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

var serviceTemplate string
var validServiceName = regexp.MustCompile(`^[a-zA-Z\-]+$`).MatchString

var terraformModuleInit = `
module "service-{{.ServiceName}}" {
	source                  = "../services/{{.ServiceName}}"
	container_project       = var.container_project
	environment             = var.environment
	project_id              = google_project.project.project_id
	project_number          = google_project.project.number
	region                  = var.region
	private_network         = google_compute_network.private_network.self_link
	service_depends_on      = [google_service_networking_connection.private_vpc_connection, google_project_iam_member.container_access, google_project_iam_binding.cloudsql]
}
`

type serviceTemplateVars struct {
	ServiceName string
	ProjectDir  string
	ServicesDir string
}

func initServiceTerraform(name string) error {
	projectsConfig, err := readProjectsConfig()
	if err != nil {
		return err
	}
	projectDir := projectsConfig.Projects[projectsConfig.Selected]
	serviceModulesPath := path.Join(projectDir, "infra", "services.tf")
	fmt.Println("Adding service to modules:", serviceModulesPath)

	f, err := os.OpenFile(serviceModulesPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	t, err := template.New("service-terraform-init").Parse(terraformModuleInit)
	if err != nil {
		return err
	}
	if t.Execute(f, &serviceTemplateVars{ServiceName: name}); err != nil {
		return err
	}
	return nil
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
		"ToCamel": strcase.ToCamel,
		"ToKebab": strcase.ToKebab,
		"ToSnake": strcase.ToSnake,
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

		name := args[0]
		if validServiceName(name) != true {
			return errors.New("Service names may only contain letters, dashes, and underscores")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		if err := createService(name); err != nil {
			fmt.Printf("Error creating service: %v\n", err)
			os.Exit(1)
		}
		if err := initServiceTerraform(name); err != nil {
			fmt.Printf("Error generating proto: %v\n", err)
			os.Exit(1)
		}
		servicePath := path.Join("services", name)
		if err := genProto(servicePath); err != nil {
			fmt.Printf("Error generating proto: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	serviceCmd.AddCommand(createServiceCmd)
	createServiceCmd.Flags().StringVarP(&serviceTemplate, "template", "t", "", "Template to generate service from")
}
