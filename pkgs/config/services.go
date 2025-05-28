package config

import (
	"fmt"
	"log"
	"motion/pkgs/models"

	"github.com/spf13/viper"
)

func InitServicesConfig() {
	servicesConfig := viper.New()
	
	servicesConfig.SetConfigName("services")
	servicesConfig.SetConfigType("yaml")
	servicesConfig.AddConfigPath(".")
	
	if err := servicesConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Arquivo services.yaml não encontrado. Criando novo arquivo...")
			
			servicesConfig.Set("repos", []models.RepoConfig{})
			
			err := servicesConfig.SafeWriteConfigAs("services.yaml")
			if err != nil {
				log.Fatalf("Erro ao criar services.yaml: %v", err)
			}
		} else {
			log.Fatalf("Erro ao ler services.yaml: %v", err)
		}
	}
	
	if err := servicesConfig.UnmarshalKey("repos", &Repos); err != nil {
		log.Fatalf("Erro ao carregar repositórios: %v", err)
	}
}

func LoadRepos() {
	servicesConfig := viper.New()
	servicesConfig.SetConfigName("services")
	servicesConfig.SetConfigType("yaml")
	servicesConfig.AddConfigPath(".")
	
	if err := servicesConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			Repos = make([]models.RepoConfig, 0)
			return
		}
		log.Fatalf("Erro ao ler services.yaml: %v", err)
	}
	
	if err := servicesConfig.UnmarshalKey("repos", &Repos); err != nil {
		log.Fatalf("Erro ao carregar repositórios: %v", err)
	}
}

func SaveRepos() {
	servicesConfig := viper.New()
	servicesConfig.SetConfigName("services")
	servicesConfig.SetConfigType("yaml")
	servicesConfig.AddConfigPath(".")
	
	_ = servicesConfig.ReadInConfig()
	
	servicesConfig.Set("repos", Repos)
	
	if err := servicesConfig.WriteConfig(); err != nil {
		if err := servicesConfig.SafeWriteConfigAs("services.yaml"); err != nil {
			log.Fatalf("Erro ao salvar services.yaml: %v", err)
		}
	}
}