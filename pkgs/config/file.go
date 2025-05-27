package config

import (
	"fmt"
	"log"
	"motion/core"
	"motion/pkgs/gitclient"

	"github.com/google/go-github/v55/github"
	"github.com/spf13/viper"
)

type RepoConfig struct {
	Name   string   `json:"name"`
	Branch string   `json:"branch"`
	Path   string   `json:"path"`
	Ports  [2]int   `json:"ports"`
	Events []string `json:"events"`
}

type Config struct {
	Secret      string                `mapstructure:"secret"`
	CurrentPort string                `mapstructure:"current_port"`
	Repos       map[string]RepoConfig `mapstructure:"repos"`
	GhToken     string                `mapstructure:"github_token"`
	UserName    string                `mapstructure:"username"`
	GhClient    *github.Client        `mapstructure:"-"`
}

var All Config
var Engine *core.Instance
var Repos []RepoConfig

func Init() {
	// Configuração principal (motion.yaml)
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

	// Carrega os repositórios do services.yaml
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

func AddRepo(repo RepoConfig) {
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