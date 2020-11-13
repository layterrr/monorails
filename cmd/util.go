package cmd

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func getProjectDirectory() (string, error) {
	projectsConfig, err := readProjectsConfig()
	if err != nil {
		return "", err
	}
	return projectsConfig.Projects[projectsConfig.Selected], nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func randomString(n int) string {
	letters := "abcdefghijklmnopqrstuvwxyz1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[seededRand.Intn(len(letters))]
	}
	return string(b)
}

func printCmd(description string, cmd *exec.Cmd) {
	fmt.Println("Command:", description)
	for _, v := range cmd.Env {
		fmt.Printf("%s \\\n", v)
	}
	fmt.Printf("%s", cmd)
	fmt.Println("")
}

func createCommand(command, dir string) *exec.Cmd {
	commandParts := strings.Split(command, " ")
	commandName := commandParts[0]
	commandArgs := commandParts[1:]
	for i, arg := range commandArgs {
		gopath := os.ExpandEnv("$GOPATH")
		arg = strings.ReplaceAll(arg, "$GOPATH", gopath)
		arg = strings.ReplaceAll(arg, "$CURRENT_MONORAILS_DIR", "projectDir")
		commandArgs[i] = arg
	}
	fmt.Println(commandArgs)
	cmd := exec.Command(commandName, commandArgs...)
	cmd.Dir = dir
	return cmd
}

func runCommand(description string, cmd *exec.Cmd) (string, error) {
	if cmd.Dir == "" {
		cmd.Dir = "projectDir"
	}
	printCmd(description, cmd)
	var out []byte
	var err error
	if !dryRun {
		out, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Output:", string(out))
			fmt.Println("Error:", err)
			return "", err
		}
		fmt.Println("Output:", string(out))
	}

	return string(out), nil
}

func echoCommand(description string, cmd *exec.Cmd) (string, error) {
	fmt.Println(cmd)
	return "", nil
}

func cloudrunServiceAccount(projectNumber string) string {
	return fmt.Sprintf("serviceAccount:service-%s@serverless-robot-prod.iam.gserviceaccount.com", projectNumber)
}

func defaultServiceAccount(projectNumber string) string {
	return fmt.Sprintf("serviceAccount:%s-compute@developer.gserviceaccount.com", projectNumber)
}

func findAndParseTemplates(rootDir string) (*template.Template, error) {
	cleanRoot := filepath.Clean(rootDir)
	pfx := len(cleanRoot) + 1
	root := template.New("")

	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			if e1 != nil {
				return e1
			}

			b, e2 := ioutil.ReadFile(path)
			if e2 != nil {
				return e2
			}

			name := path[pfx:]
			t := root.New(name)
			_, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}
		}

		return nil
	})

	return root, err
}

func copyTemplate(dir string, data interface{}) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, e1 error) error {
		fmt.Println(path)
		return nil
	})
	return err
}
