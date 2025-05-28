package models

import "github.com/google/go-github/v55/github"

type RepoConfig struct {
	GitID      string   `json:"git_id"`
	Name    string   `json:"name"`
	Branch  string   `json:"branch"`
	Path    string   `json:"path"`
	Ports   [2]int   `json:"ports"`
	Events  []string `json:"events"`
	Service CoreApplication
}

type Config struct {
	Secret      string                `mapstructure:"secret"`
	CurrentPort string                `mapstructure:"current_port"`
	Repos       map[string]RepoConfig `mapstructure:"repos"`
	GhToken     string                `mapstructure:"github_token"`
	UserName    string                `mapstructure:"username"`
	GhClient    *github.Client        `mapstructure:"-"`
}