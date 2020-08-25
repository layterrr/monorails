package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var force bool
var dryRun bool
var stagingProject = getEnv("MONORAILS_STAGING_PROJECT_ID", "")
var productionProject = getEnv("MONORAILS_PRODUCTION_PROJECT_ID", "")
var testProject = getEnv("MONORAILS_TEST_PROJECT_ID", stagingProject)
var containerProjectID = getEnv("MONORAILS_CONTAINER_PROJECT_ID", stagingProject)
var targetEnvironment string
var isProd bool
var isStaging bool
var isTest bool
var environments = map[string]string{
	"staging":    stagingProject,
	"production": productionProject,
	"test":       testProject,
}

var rootCmd = &cobra.Command{
	Use:   "monorails",
	Short: "Utilities for managing a monorepo",
}

// Execute run root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "n", false, "log any commands with side effects")
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Force command. Required to run dangerous actions on protected projects.")
	rootCmd.PersistentFlags().BoolVar(&isProd, "prod", false, "Run commands against production")
	rootCmd.PersistentFlags().BoolVar(&isStaging, "staging", false, "Run commands against staging")
	rootCmd.PersistentFlags().BoolVar(&isTest, "test", false, "Run commands against test")
}
