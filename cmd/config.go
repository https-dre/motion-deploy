package cmd

import (
	"motion/pkgs/config"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	username     string
	gitToken     string
	port         string
	serverSecret string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Define configurações específicas do sistema",
	Run: func(cmd *cobra.Command, args []string) {
		if username != "" {
			config.All.UserName = username
			fmt.Println("Username configurado para:", username)
		}
		if gitToken != "" {
			config.All.GhToken = gitToken
			fmt.Println("Token do Git configurado.")
		}
		if port != "" {
			config.All.CurrentPort = port
			fmt.Println("Porta configurada para:", port)
		}
		if serverSecret != "" {
			config.All.Secret = serverSecret
			fmt.Println("Segredo do servidor configurado.")
		}

		if username == "" && gitToken == "" && port == "" && serverSecret == "" {
			fmt.Println("Nenhuma opção fornecida. Use --help para ver as opções disponíveis.")
		}

		config.Save()
	},
}

func init() {
	configCmd.Flags().StringVar(&username, "username", "", "Nome de usuário para deploy")
	configCmd.Flags().StringVar(&gitToken, "git_token", "", "Token de autenticação do Git")
	configCmd.Flags().StringVar(&port, "port", "", "Porta do servidor HTTP")
	configCmd.Flags().StringVar(&serverSecret, "secret", "", "Segredo do servidor para autenticação")
}
