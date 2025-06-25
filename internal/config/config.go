package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	GitHubToken     string
	AnthropicAPIKey string
	ClaudeSkipPerm  bool
	CacheDir        string
	DockerImage     string
	DefaultBranch   string
	SetupScript     string
}

func Load() (*Config, error) {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable is required")
	}

	anthropicAPIKey := os.Getenv("ANTHROPIC_API_KEY")
	if anthropicAPIKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable is required")
	}

	cacheDir := os.Getenv("CACHE_DIR")
	if cacheDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		cacheDir = filepath.Join(homeDir, ".cache", "claudecli")
	}

	claudeSkipPerm := os.Getenv("CLAUDE_SKIP_PERM") != ""
	setupScript := os.Getenv("SETUP_SCRIPT")

	return &Config{
		GitHubToken:     githubToken,
		AnthropicAPIKey: anthropicAPIKey,
		ClaudeSkipPerm:  claudeSkipPerm,
		CacheDir:        cacheDir,
		DockerImage:     "ghcr.io/anthropic/claude-code-devcontainer:latest",
		DefaultBranch:   "main",
		SetupScript:     setupScript,
	}, nil
}
