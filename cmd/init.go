package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Inicia configuração",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Iniciando configuração!")
	},
}
