package cmd

import (
	"fmt"
	"log"
	"motion/pkgs/config"
	"motion/pkgs/repo"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	ports  string
	branch string
)

var add = &cobra.Command{
	Use:   "add <reponame> [flags]",
	Short: "Adiciona um repositório",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		reponame := args[0]
		ports_selected := [2]int{0, 2}
		repository, err := repo.FindRepo(reponame)

		if err != nil {
			log.Fatal(err)
		}

		if repository == nil {
			fmt.Println("REPOSITORY NOT FOUND!")
			os.Exit(1)
		}

		branch_selected := repository.DefaultBranch

		if branch != "" {
			branch_selected = &branch
		}

		if ports != "" {
			arr := strings.Split(ports, ":")

			if len(arr) != 2 {
				fmt.Print("Formato inválido. Esperado: porta1:porta2")
			}

			port1, err1 := strconv.Atoi(arr[0])
			port2, err2 := strconv.Atoi(arr[1])

			if err1 != nil || err2 != nil {
				fmt.Print("Erro ao converter portas para inteiro.")
			}

			ports_selected[0] = port1
			ports_selected[1] = port2
		}

		path := filepath.Join("./services", reponame)

		fmt.Printf("Adding Repository: %s/%s\n", config.General.UserName, reponame)

		new_repo := config.RepoConfig{
			Name:   reponame,
			Ports:  ports_selected,
			Branch: *branch_selected,
			Events: []string{"push"},
			Path:   path,
		}

		repo.DownloadRepository(repository, "./services")
		config.General.AddRepo(new_repo)
		config.General.Save()
		fmt.Println("Repository downloaded!")
		fmt.Println("Setting project in Docker")

		if err := config.Engine.BuildAndRunService(new_repo.Path, new_repo.Name, new_repo.Name); err != nil {
			fmt.Println("Falha ao buildar e rodar o projeto")
			fmt.Println(err)
		}
	},
}

func init() {
	add.Flags().StringVar(&ports, "ports", "", "Define port mapping")
	add.Flags().StringVar(&branch, "branch", "", "Branch selected")
}
