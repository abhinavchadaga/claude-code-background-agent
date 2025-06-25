package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v57/github"
)

type Service struct {
	client *github.Client
	ctx    context.Context
}

type PRConfig struct {
	Owner       string
	Repo        string
	Title       string
	Body        string
	Head        string
	Base        string
	IssueNumber int
	Labels      []string
}

func NewService(client *github.Client, ctx context.Context) *Service {
	return &Service{
		client: client,
		ctx:    ctx,
	}
}

func (s *Service) CreatePullRequest(config PRConfig) (*github.PullRequest, error) {
	pullRequest := &github.NewPullRequest{
		Title: &config.Title,
		Body:  &config.Body,
		Head:  &config.Head,
		Base:  &config.Base,
	}

	pr, _, err := s.client.PullRequests.Create(s.ctx, config.Owner, config.Repo, pullRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create pull request: %w", err)
	}

	if len(config.Labels) > 0 {
		_, _, err = s.client.Issues.AddLabelsToIssue(s.ctx, config.Owner, config.Repo, pr.GetNumber(), config.Labels)
		if err != nil {
			fmt.Printf("Warning: failed to add labels to PR: %v\n", err)
		}
	}

	if config.IssueNumber > 0 {
		comment := fmt.Sprintf("This pull request addresses issue #%d", config.IssueNumber)
		issueComment := &github.IssueComment{
			Body: &comment,
		}
		_, _, err = s.client.Issues.CreateComment(s.ctx, config.Owner, config.Repo, config.IssueNumber, issueComment)
		if err != nil {
			fmt.Printf("Warning: failed to link PR to issue: %v\n", err)
		}
	}

	return pr, nil
}

func (s *Service) BranchExists(owner, repo, branch string) (bool, error) {
	_, _, err := s.client.Git.GetRef(s.ctx, owner, repo, fmt.Sprintf("heads/%s", branch))
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Service) CreateBranch(owner, repo, newBranch, baseBranch string) error {
	baseRef, _, err := s.client.Git.GetRef(s.ctx, owner, repo, fmt.Sprintf("heads/%s", baseBranch))
	if err != nil {
		return fmt.Errorf("failed to get base branch reference: %w", err)
	}

	newRef := &github.Reference{
		Ref:    github.String(fmt.Sprintf("refs/heads/%s", newBranch)),
		Object: baseRef.Object,
	}

	_, _, err = s.client.Git.CreateRef(s.ctx, owner, repo, newRef)
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	return nil
}

func (s *Service) GetDefaultBranch(owner, repo string) (string, error) {
	repository, _, err := s.client.Repositories.Get(s.ctx, owner, repo)
	if err != nil {
		return "", fmt.Errorf("failed to get repository info: %w", err)
	}
	return repository.GetDefaultBranch(), nil
}

func GeneratePRTitle(issueTitle string, issueNumber int) string {
	return fmt.Sprintf("Fix #%d: %s", issueNumber, issueTitle)
}

func GeneratePRBody(issueNumber int, issueTitle, issueBody string) string {
	var body strings.Builder

	body.WriteString(fmt.Sprintf("## 🔧 Automated Fix for Issue #%d\n\n", issueNumber))
	body.WriteString(fmt.Sprintf("**Original Issue:** %s\n\n", issueTitle))

	if issueBody != "" {
		body.WriteString("### Issue Description\n")

		lines := strings.Split(issueBody, "\n")
		for _, line := range lines {
			body.WriteString("> ")
			body.WriteString(line)
			body.WriteString("\n")
		}
		body.WriteString("\n")
	}

	body.WriteString("### Changes Made\n")
	body.WriteString("This pull request was generated automatically by Claude Code to address the issue described above. ")
	body.WriteString("The changes have been analyzed and implemented to resolve the reported problem.\n\n")

	body.WriteString("### Review Notes\n")
	body.WriteString("- 🤖 This PR was created by Claude Code Background Agent\n")
	body.WriteString("- 📋 Please review the changes carefully before merging\n")
	body.WriteString("- 🧪 Test the changes in your development environment\n")
	body.WriteString("- 🔍 Verify that the original issue is resolved\n\n")

	body.WriteString(fmt.Sprintf("Closes #%d", issueNumber))

	return body.String()
}

func GetSuggestedLabels(issue *github.Issue) []string {
	labels := []string{"automated-fix", "claude-code"}

	for _, label := range issue.Labels {
		labelName := label.GetName()
		if strings.Contains(labelName, "bug") {
			labels = append(labels, "bug-fix")
		} else if strings.Contains(labelName, "enhancement") || strings.Contains(labelName, "feature") {
			labels = append(labels, "enhancement")
		} else if strings.Contains(labelName, "documentation") {
			labels = append(labels, "documentation")
		}
	}

	return labels
}
