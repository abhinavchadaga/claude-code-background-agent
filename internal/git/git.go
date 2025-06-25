package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type RepoInfo struct {
	Owner      string
	Name       string
	RemoteURL  string
	RootPath   string
	CurrentSHA string
}

func GetCurrentRepo() (*RepoInfo, error) {
	rootPath, err := getGitRoot()
	if err != nil {
		return nil, fmt.Errorf("not in a git repository: %w", err)
	}

	remoteURL, err := getRemoteURL(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get remote URL: %w", err)
	}

	owner, name, err := parseGitHubURL(remoteURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GitHub URL: %w", err)
	}

	currentSHA, err := getCurrentSHA(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get current SHA: %w", err)
	}

	return &RepoInfo{
		Owner:      owner,
		Name:       name,
		RemoteURL:  remoteURL,
		RootPath:   rootPath,
		CurrentSHA: currentSHA,
	}, nil
}

func getGitRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getRemoteURL(repoPath string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getCurrentSHA(repoPath string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func parseGitHubURL(url string) (owner, name string, err error) {
	url = strings.TrimSuffix(url, ".git")

	if strings.HasPrefix(url, "git@github.com:") {
		parts := strings.Split(strings.TrimPrefix(url, "git@github.com:"), "/")
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid SSH GitHub URL format")
		}
		return parts[0], parts[1], nil
	}

	if strings.HasPrefix(url, "https://github.com/") {
		parts := strings.Split(strings.TrimPrefix(url, "https://github.com/"), "/")
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid HTTPS GitHub URL format")
		}
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("not a GitHub repository URL")
}

func CloneRepo(url, targetDir string, shallow bool) error {
	if err := os.MkdirAll(filepath.Dir(targetDir), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	args := []string{"clone"}
	if shallow {
		args = append(args, "--depth=1")
	}
	args = append(args, url, targetDir)

	cmd := exec.Command("git", args...)
	return cmd.Run()
}

func CreateCacheDir(cacheDir, owner, repo, sha string) string {
	return filepath.Join(cacheDir, "repos", fmt.Sprintf("%s-%s", owner, repo), sha)
}

func CloneToCache(repoInfo *RepoInfo, cacheDir string) (string, error) {
	targetDir := CreateCacheDir(cacheDir, repoInfo.Owner, repoInfo.Name, repoInfo.CurrentSHA)

	if _, err := os.Stat(targetDir); err == nil {
		return targetDir, nil
	}

	fmt.Printf("Cloning repository to cache: %s\n", targetDir)
	if err := CloneRepo(repoInfo.RemoteURL, targetDir, true); err != nil {
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}

	return targetDir, nil
}

func HasChanges(repoPath string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %w", err)
	}

	return len(strings.TrimSpace(string(output))) > 0, nil
}

func CreateBranch(repoPath, branchName string) error {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create branch %s: %w", branchName, err)
	}
	return nil
}

func CommitChanges(repoPath, message string) error {
	addCmd := exec.Command("git", "add", ".")
	addCmd.Dir = repoPath
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Dir = repoPath
	if err := commitCmd.Run(); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	return nil
}

func PushBranch(repoPath, branchName, remoteURL string) error {
	cmd := exec.Command("git", "push", "origin", branchName)
	cmd.Dir = repoPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push branch %s: %w", branchName, err)
	}

	return nil
}

func GetCommitMessage(repoPath string) (string, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=format:%s")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get last commit message: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func ConfigureGit(repoPath, name, email string) error {
	nameCmd := exec.Command("git", "config", "user.name", name)
	nameCmd.Dir = repoPath
	if err := nameCmd.Run(); err != nil {
		return fmt.Errorf("failed to configure git user name: %w", err)
	}

	emailCmd := exec.Command("git", "config", "user.email", email)
	emailCmd.Dir = repoPath
	if err := emailCmd.Run(); err != nil {
		return fmt.Errorf("failed to configure git user email: %w", err)
	}

	return nil
}
