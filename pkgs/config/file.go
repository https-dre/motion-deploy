package config

import (
	"fmt"
	"log"
	"motion/pkgs/gitclient"
	"motion/pkgs/models"

	"github.com/spf13/viper"
)



var All models.Config
var Repos []models.RepoConfig

func Init() {
	viper.SetConfigName("motion")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("current_port", "5500")
	viper.SetDefault("secret", "Undefined")
	viper.SetDefault("username", "Undefined")
	viper.SetDefault("github_token", "Undefined")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Arquivo de Configuração não encontrado. Criando motion.yaml")
			err := viper.SafeWriteConfigAs("motion.yaml")
			if err != nil {
				log.Fatalf("Erro ao criar motion.yaml: %v", err)
			}
		} else {
			log.Fatalf("Erro ao ler motion.yaml: %v", err)
		}
	}

	if err := viper.Unmarshal(&All); err != nil {
		log.Fatalf("Erro ao carregar config para struct: %v", err)
	}

	LoadRepos()
}

func InitGitClient() {
	if All.GhToken == "" || All.GhToken == "Undefined" {
		fmt.Println("Invalid GitHub token!")
		return
	}
	All.GhClient = gitclient.NewGitClient(All.GhToken)
}

func Save() error {
	viper.Set("current_port", All.CurrentPort)
	viper.Set("secret", All.Secret)
	viper.Set("username", All.UserName)
	viper.Set("github_token", All.GhToken)

	return viper.WriteConfig()
}

func AddRepo(repo models.RepoConfig) {
	Repos = append(Repos, repo)
	SaveRepos()
}

func RemoveRepo(repoName string) {
	for i, repo := range Repos {
		if repo.Name == repoName {
			Repos = append(Repos[:i], Repos[i+1:]...)
			break
		}
	}
	SaveRepos()
}
