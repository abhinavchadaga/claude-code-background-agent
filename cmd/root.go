package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	logLevel string
)

var rootCmd = &cobra.Command{
	Use:   "claudecli",
	Short: "Claude Code Background Agent CLI",
	Long: `A command-line tool that automates GitHub issue resolution using Claude Code
in isolated Docker containers. Handles clarification, planning, and implementation.`,
}

var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "Commands for managing Docker containers",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level (debug|info|warn|error)")
	rootCmd.AddCommand(containerCmd)
}
