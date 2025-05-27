package cmd

import (
	"context"
	"fmt"
	"motion/pkgs/config"
	"net/http"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/google/go-github/v55/github"
	"golang.org/x/oauth2"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func input(prompt string) string {
	var value string
	survey.AskOne(&survey.Input{Message: prompt}, &value)
	return value
}

func githubProfileExists(username string) (bool, error) {
	client := config.All.GhClient

	_, resp, err := client.Users.Get(context.Background(), username)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func verifyGitToken(token string) (bool, error){
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	_, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Config Init",
	Run: func(cmd *cobra.Command, args []string) {
		color.HiGreen("Thank you for choose the motion-deploy!")

		username := input("Your GitHub username>")
		if username == "" {
			color.Red("Please provide an username")
			os.Exit(0)
		}

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // ⠋ charset animado

		s.Suffix = " Verifying your profile..."
		s.Start()
		result, err := githubProfileExists(username)

		if err != nil {
			s.Stop()
			color.Red("An error ocurred in profile verification...")
			fmt.Println("Try again!")
			os.Exit(1)
		}

		if !result {
			s.Stop()
			color.Red("GitHub profile not found!")
			os.Exit(0)
		}

		time.Sleep(1 * time.Second)
		s.Stop()
		fmt.Println("Profile ok!")

		fmt.Printf("\nNow, we need add one %s...\n", color.RedString("GitHub Token"))
		fmt.Println("In your GitHub, go to Settings > Developer Settings > Tokens (classic)")
		fmt.Println("Generate an new Token with no Expiration")
		fmt.Println("Mark all 'repo' and 'admin:repo_hook' scope")
		fmt.Print("Add an note to your token\n\n")

		var token string
		survey.AskOne(&survey.Password{
			Message: "Past your token here>",
		}, &token)

		s.Suffix = " Checking your GitHub token..."
		s.Start()
		token_result, err := verifyGitToken(token)
		if err != nil {
			s.Stop()
			color.Red("An error ocurred in token verification!")
			os.Exit(1)
		}

		if !token_result {
			s.Stop()
			color.Red("Token invalid!")
			os.Exit(0)
		}

		time.Sleep(1 * time.Second)
		s.Stop()
		fmt.Println("Token verified!")

		port := input("Choose one port for motion to run [default: 5500]>")

		config.All.UserName = username
		config.All.GhToken = token
		config.All.CurrentPort = port
		config.InitServicesConfig()
		config.Save()

		color.Blue("\nMotion configured!")
		fmt.Println("All motion configs was saved in 'motion.yaml' file.")
	},
}
