package claude

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v57/github"
)

type TaskConfig struct {
	Issue      *github.Issue
	RepoOwner  string
	RepoName   string
	BranchName string
}

func GenerateTaskPrompt(config TaskConfig) string {
	var prompt strings.Builder

	prompt.WriteString("You are Claude Code, an AI assistant that helps resolve GitHub issues by analyzing code and making changes.\n\n")

	prompt.WriteString("## Issue Details\n")
	prompt.WriteString(fmt.Sprintf("**Repository:** %s/%s\n", config.RepoOwner, config.RepoName))
	prompt.WriteString(fmt.Sprintf("**Issue #%d:** %s\n", config.Issue.GetNumber(), config.Issue.GetTitle()))
	prompt.WriteString(fmt.Sprintf("**Branch:** %s\n\n", config.BranchName))

	if body := config.Issue.GetBody(); body != "" {
		prompt.WriteString("**Issue Description:**\n")
		prompt.WriteString(body)
		prompt.WriteString("\n\n")
	}

	if labels := config.Issue.Labels; len(labels) > 0 {
		prompt.WriteString("**Labels:** ")
		for i, label := range labels {
			if i > 0 {
				prompt.WriteString(", ")
			}
			prompt.WriteString(label.GetName())
		}
		prompt.WriteString("\n\n")
	}

	prompt.WriteString("## Task\n")
	prompt.WriteString("Please analyze this repository and work on resolving the GitHub issue described above. ")
	prompt.WriteString("You have full access to the codebase and can:\n")
	prompt.WriteString("- Examine the current code structure\n")
	prompt.WriteString("- Identify the root cause of the issue\n")
	prompt.WriteString("- Make necessary code changes\n")
	prompt.WriteString("- Create, modify, or delete files as needed\n")
	prompt.WriteString("- Test your changes if possible\n\n")

	prompt.WriteString("Please start by exploring the repository to understand its structure and then work on addressing the issue. ")
	prompt.WriteString("If you need clarification or encounter any blockers, please explain your findings and next steps.")

	return prompt.String()
}

func BuildClaudeCommand(apiKey string, skipPermissions bool, taskPrompt string) []string {
	cmd := []string{"claude"}

	if skipPermissions {
		cmd = append(cmd, "--dangerously-skip-permissions")
	}

	cmd = append(cmd, "--")
	cmd = append(cmd, taskPrompt)

	return cmd
}

func BuildClaudeEnvironment(apiKey string, issueNumber int, issueTitle, repoOwner, repoName string) []string {
	env := []string{
		fmt.Sprintf("ANTHROPIC_API_KEY=%s", apiKey),
		fmt.Sprintf("ISSUE_NUMBER=%d", issueNumber),
		fmt.Sprintf("ISSUE_TITLE=%s", issueTitle),
		fmt.Sprintf("REPO_OWNER=%s", repoOwner),
		fmt.Sprintf("REPO_NAME=%s", repoName),
	}
	return env
}
