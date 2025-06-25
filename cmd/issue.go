package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/abhinavchadaga/claude-code-background-agent/internal/claude"
	"github.com/abhinavchadaga/claude-code-background-agent/internal/config"
	"github.com/abhinavchadaga/claude-code-background-agent/internal/docker"
	"github.com/abhinavchadaga/claude-code-background-agent/internal/git"
	githubservice "github.com/abhinavchadaga/claude-code-background-agent/internal/github"
	"github.com/google/go-github/v57/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var (
	issueNumber int
	branchName  string
	askMode     string
	setupScript string
	createPR    bool
	gitName     string
	gitEmail    string
)

var issueCmd = &cobra.Command{
	Use:   "issue [issue-number]",
	Short: "Process a GitHub issue and create a pull request with Claude Code",
	Long: `Downloads the specified GitHub issue, clones the repository to a cached location,
runs Claude Code in a Docker container to analyze and fix the issue, then creates a pull request with the changes.`,
	Args: cobra.ExactArgs(1),
	RunE: runIssue,
}

func init() {
	rootCmd.AddCommand(issueCmd)
	issueCmd.Flags().IntVar(&issueNumber, "number", 0, "issue number (required)")
	issueCmd.Flags().StringVar(&branchName, "branch", "", "working branch name (default: issue-<num>)")
	issueCmd.Flags().StringVar(&askMode, "ask-mode", "auto", "clarification mode (auto|interactive|none)")
	issueCmd.Flags().StringVar(&setupScript, "setup-script", "", "path to setup script to run in container")
	issueCmd.Flags().BoolVar(&createPR, "create-pr", true, "Create a pull request after successful execution")
	issueCmd.Flags().StringVar(&gitName, "git-name", "Claude Code Bot", "Git user name for commits")
	issueCmd.Flags().StringVar(&gitEmail, "git-email", "claude-code-bot@example.com", "Git user email for commits")

	issueCmd.MarkFlagRequired("number")
}

func runIssue(cmd *cobra.Command, args []string) error {
	issueNumber, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid issue number: %w", err)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	repo, err := git.GetCurrentRepo()
	if err != nil {
		return fmt.Errorf("failed to get current repository: %w", err)
	}

	if branchName == "" {
		branchName = fmt.Sprintf("issue-%d", issueNumber)
	}

	fmt.Printf("Processing issue #%d for %s/%s\n", issueNumber, repo.Owner, repo.Name)
	fmt.Printf("Using branch: %s\n", branchName)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitHubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	issue, _, err := client.Issues.Get(ctx, repo.Owner, repo.Name, issueNumber)
	if err != nil {
		return fmt.Errorf("failed to fetch issue: %w", err)
	}

	fmt.Printf("Issue title: %s\n", issue.GetTitle())
	fmt.Printf("Issue state: %s\n", issue.GetState())

	if issue.GetState() != "open" {
		return fmt.Errorf("issue #%d is not open", issueNumber)
	}

	fmt.Println("\n--- Issue Body ---")
	fmt.Println(issue.GetBody())
	fmt.Println("--- End Issue Body ---\n")

	fmt.Printf("Repository: %s\n", repo.RootPath)
	fmt.Printf("Current SHA: %s\n", repo.CurrentSHA)

	fmt.Println("Step 1: Cloning repository to cache...")
	cachedRepoPath, err := git.CloneToCache(repo, cfg.CacheDir)
	if err != nil {
		return fmt.Errorf("failed to clone repository to cache: %w", err)
	}
	fmt.Printf("Repository cached at: %s\n", cachedRepoPath)

	fmt.Println("\nStep 2: Generating Claude Code task...")
	taskConfig := claude.TaskConfig{
		Issue:      issue,
		RepoOwner:  repo.Owner,
		RepoName:   repo.Name,
		BranchName: branchName,
	}
	taskPrompt := claude.GenerateTaskPrompt(taskConfig)

	fmt.Println("Generated task prompt:")
	fmt.Println("--- Task Prompt ---")
	fmt.Println(taskPrompt)
	fmt.Println("--- End Task Prompt ---\n")

	fmt.Println("Step 3: Starting Docker container with Claude Code...")
	dockerService, err := docker.NewService()
	if err != nil {
		return fmt.Errorf("failed to create Docker service: %w", err)
	}
	defer dockerService.Close()

	scriptPath := cfg.SetupScript
	if setupScript != "" {
		scriptPath = setupScript
	}

	claudeCommand := claude.BuildClaudeCommand(cfg.AnthropicAPIKey, cfg.ClaudeSkipPerm, taskPrompt)
	claudeEnv := claude.BuildClaudeEnvironment(cfg.AnthropicAPIKey, issueNumber, issue.GetTitle(), repo.Owner, repo.Name)

	containerConfig := docker.ContainerConfig{
		Image:       cfg.DockerImage,
		RepoPath:    cachedRepoPath,
		WorkDir:     "/workspace",
		Command:     claudeCommand,
		Env:         claudeEnv,
		SetupScript: scriptPath,
	}

	containerID, err := dockerService.StartContainer(ctx, containerConfig)
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	fmt.Printf("Container started with ID: %s\n", containerID[:12])
	fmt.Println("Claude Code is now analyzing the repository and working on the issue...")

	fmt.Println("\nStep 4: Streaming Claude Code logs...")
	logs, err := dockerService.GetContainerLogs(ctx, containerID)
	if err != nil {
		return fmt.Errorf("failed to get container logs: %w", err)
	}
	defer logs.Close()

	fmt.Println("--- Claude Code Output ---")
	_, err = io.Copy(os.Stdout, logs)
	if err != nil {
		return fmt.Errorf("failed to stream logs: %w", err)
	}
	fmt.Println("--- End Claude Code Output ---")

	fmt.Println("\nStep 5: Waiting for Claude Code to complete...")
	if err := dockerService.WaitForContainer(ctx, containerID); err != nil {
		return fmt.Errorf("Claude Code execution failed: %w", err)
	}

	fmt.Printf("\n🎉 Milestone 3 complete! Claude Code executed successfully in container %s.\n", containerID[:12])
	fmt.Println("\nImplemented features:")
	fmt.Println("✅ Claude Code invocation with task generation")
	fmt.Println("✅ Real-time log streaming from Claude execution")
	fmt.Println("✅ Proper exit code handling")
	fmt.Println("✅ GitHub issue analysis and processing")
	fmt.Println("✅ Development environment setup integration")

	if !createPR {
		fmt.Printf("Skipping PR creation (--create-pr=false)\n")
		return nil
	}

	fmt.Printf("Checking for repository changes...\n")
	hasChanges, err := git.HasChanges(cachedRepoPath)
	if err != nil {
		return fmt.Errorf("failed to check for changes: %w", err)
	}

	if !hasChanges {
		fmt.Printf("No changes detected in repository, skipping PR creation\n")
		return nil
	}

	fmt.Printf("Changes detected, proceeding with PR creation...\n")

	if err := git.ConfigureGit(cachedRepoPath, gitName, gitEmail); err != nil {
		return fmt.Errorf("failed to configure git: %w", err)
	}

	branchName := fmt.Sprintf("claude-code/fix-issue-%d", issueNumber)
	fmt.Printf("Creating branch: %s\n", branchName)
	if err := git.CreateBranch(cachedRepoPath, branchName); err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	commitMessage := fmt.Sprintf("Fix #%d: %s\n\nAutomatically generated fix by Claude Code", issueNumber, issue.GetTitle())
	fmt.Printf("Committing changes...\n")
	if err := git.CommitChanges(cachedRepoPath, commitMessage); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	fmt.Printf("Pushing branch to GitHub...\n")
	if err := git.PushBranch(cachedRepoPath, branchName, repo.RemoteURL); err != nil {
		return fmt.Errorf("failed to push branch: %w", err)
	}

	githubSvc := githubservice.NewService(client, ctx)

	defaultBranch, err := githubSvc.GetDefaultBranch(repo.Owner, repo.Name)
	if err != nil {
		return fmt.Errorf("failed to get default branch: %w", err)
	}

	prConfig := githubservice.PRConfig{
		Owner:       repo.Owner,
		Repo:        repo.Name,
		Title:       githubservice.GeneratePRTitle(issue.GetTitle(), issueNumber),
		Body:        githubservice.GeneratePRBody(issueNumber, issue.GetTitle(), issue.GetBody()),
		Head:        branchName,
		Base:        defaultBranch,
		IssueNumber: issueNumber,
		Labels:      githubservice.GetSuggestedLabels(issue),
	}

	fmt.Printf("Creating pull request...\n")
	pr, err := githubSvc.CreatePullRequest(prConfig)
	if err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}

	fmt.Printf("✅ Successfully created pull request: %s\n", pr.GetHTMLURL())
	fmt.Printf("📋 PR #%d: %s\n", pr.GetNumber(), pr.GetTitle())

	return nil
}
