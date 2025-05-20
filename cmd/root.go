package cmd

import (
	"fmt"
	"os"

	"motion/pkgs/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "motion",
	Short: "Ferramenta para deploy automatizado via webhook",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.General.Init()
		config.General.InitGitClient()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(add)
}
