package cmd

import (
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:     "project",
	Short:   "Group command for managing projects",
	Aliases: []string{"projects"},
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
