package cmd

import (
	"fmt"
	"motion/core"
	"motion/pkgs/config"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "motion",
	Short: "Ferramenta para deploy automatizado via webhook",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.Init()
		config.InitGitClient()
		core.Engine = core.NewDockerCore()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(add)
}
