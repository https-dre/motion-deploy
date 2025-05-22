package repo

import (
	"motion/pkgs/config"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"github.com/google/go-github/v55/github"
)

func FindRepo(reponame string) (*github.Repository, error) {
	ctx := context.Background()
	client := config.General.GhClient
	username := config.General.UserName
	repos, _, err := client.Repositories.List(ctx, username, nil)

	if err != nil {
		log.Fatalf("Erro ao listar repositórios: %v", err)
		return nil, err
	}

	for _, repo := range repos {
		if reponame == repo.GetName() {
			return repo, nil
		}
	}

	return nil, nil
}

func DownloadRepository(repo *github.Repository, path string) error {
	os.MkdirAll(path, os.ModePerm)
	
	if repo != nil {
		reponame := repo.GetName()
		cloneUrl := repo.GetCloneURL()
		fmt.Printf("Clonando repositório %s em %s\n", cloneUrl, path)
		if err := cloneRepo(cloneUrl, fmt.Sprintf("%s/%s", path, reponame)); err != nil {
			log.Printf("Erro ao clonar %s: %v", reponame, err)
		}
		return nil
	}
	return nil
}

func cloneRepo(cloneURL, path string) error {
	// Verifica se a pasta já existe
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Printf("Pasta %s já existe, ignorando...\n", path)
		return nil
	}

	// Executa `git clone`
	cmd := exec.Command("git", "clone", cloneURL, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
