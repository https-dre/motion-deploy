package config

import (
	"encoding/json"
	"fmt"
	"log"
	"motion/core"
	"motion/pkgs/gitclient"
	"os"

	"github.com/google/go-github/v55/github"
)

type RepoConfig struct {
	Name   string   `json:"name"`
	Branch string   `json:"branch"`
	Path   string   `json:"path"`
	Ports  [2]int   `json:"ports"`
	Events []string `json:"events"`
}

type Config struct {
	Secret      string                `json:"secret"`
	CurrentPort string                `json:"current_port"`
	Repos       map[string]RepoConfig `json:"repos"`
	GhClient    *github.Client
	GhToken     string `json:"GITHUB_TOKEN"`
	UserName    string `json:"username"`
}

var General Config
var Engine core.Instance

func (c *Config) Init() {
	// Carrega config.json
	_, err := os.Stat("config.json")

	if err == nil { /* arquivo existe */
		data, err := os.ReadFile("config.json")
		if err != nil {
			log.Fatal("Erro ao ler config.json:", err)
		}

		if err := json.Unmarshal(data, c); err != nil {
			log.Fatal("Erro ao parsear config:", err)
		}
	}

	if os.IsNotExist(err) {
		fmt.Println("Arquivo de Configuração não existe! Criando config.json")
		c.CurrentPort = "5500"
		c.Secret = "Undefined"
		c.UserName = "Undefined"
		c.GhToken = "Undefined"
		c.Save()
	}

}

func (c *Config) InitGitClient() {
	if c.GhToken == "" {
		fmt.Println("GITHUB TOKEN INVÁLIDO")
		return
	}
	c.GhClient = gitclient.NewGitClient(General.GhToken)
}

func (c *Config) Save() error {
	config_json, err := json.MarshalIndent(c, "", "	")

	if err != nil {
		return err
	}

	err = os.WriteFile("config.json", config_json, 0644)

	if err != nil {
		return err
	}

	return nil
}

func (c *Config) AddRepo(repo RepoConfig) {
	if c.Repos == nil {
		c.Repos = make(map[string]RepoConfig)
	}

	c.Repos[repo.Name] = repo
}
