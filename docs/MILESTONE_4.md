# Milestone 4: GitHub PR Creation & Label Management

## Overview

Milestone 4 implements automatic GitHub pull request creation after successful Claude Code execution, complete with intelligent label management and issue linking.

## Features Implemented

### 1. GitHub Service (`internal/github/github.go`)

A comprehensive GitHub service that handles:

- **Pull Request Creation**: Automated PR creation with rich metadata
- **Branch Management**: Branch existence checking and creation
- **Label Management**: Intelligent label assignment based on issue type
- **Issue Linking**: Automatic linking between issues and PRs

#### Key Functions

```go
// Create pull requests with labels and issue linking
func (s *Service) CreatePullRequest(config PRConfig) (*github.PullRequest, error)

// Manage branches for PR workflow
func (s *Service) BranchExists(owner, repo, branch string) (bool, error)
func (s *Service) CreateBranch(owner, repo, newBranch, baseBranch string) error

// Get repository metadata
func (s *Service) GetDefaultBranch(owner, repo string) (string, error)

// Generate PR content
func GeneratePRTitle(issueTitle string, issueNumber int) string
func GeneratePRBody(issueNumber int, issueTitle, issueBody string) string
func GetSuggestedLabels(issue *github.Issue) []string
```

### 2. Enhanced Git Operations (`internal/git/git.go`)

Extended git operations for PR workflow:

```go
// Detect repository changes after Claude Code execution
func HasChanges(repoPath string) (bool, error)

// Git workflow functions
func CreateBranch(repoPath, branchName string) error
func CommitChanges(repoPath, message string) error
func PushBranch(repoPath, branchName, remoteURL string) error
func ConfigureGit(repoPath, name, email string) error
```

### 3. Complete Issue-to-PR Workflow (`cmd/issue.go`)

The issue command now includes full PR creation functionality:

1. **Issue Processing**: Fetch and validate GitHub issue
2. **Repository Preparation**: Clone repository to cache
3. **Claude Code Execution**: Run analysis and fixes in Docker container
4. **Change Detection**: Check if Claude Code made changes
5. **Git Workflow**: Create branch, commit changes, push to GitHub
6. **PR Creation**: Generate and create pull request with appropriate labels
7. **Issue Linking**: Automatically link PR to original issue

#### Command Options

```bash
# Process issue and create PR (default)
claudecli issue 123

# Skip PR creation
claudecli issue 123 --create-pr=false

# Custom git configuration
claudecli issue 123 --git-name="Claude Bot" --git-email="bot@example.com"
```

### 4. Intelligent Label Management

The system automatically suggests labels based on:

- **Base Labels**: Always applied `["automated-fix", "claude-code"]`
- **Issue Type Detection**:
  - Bug issues → `"bug-fix"` label
  - Enhancement/feature → `"enhancement"` label  
  - Documentation → `"documentation"` label

### 5. Rich PR Content Generation

Generated PRs include:

- **Descriptive Titles**: Format: `Fix #123: Original Issue Title`
- **Comprehensive Body**:
  - Issue summary and description
  - Automated fix explanation
  - Review guidelines
  - Automatic issue closure reference

#### Example PR Body

```markdown
## 🔧 Automated Fix for Issue #123

**Original Issue:** Fix login validation bug

### Issue Description
> User authentication fails when special characters are used in passwords.
> This prevents users from logging in with secure passwords.

### Changes Made
This pull request was generated automatically by Claude Code to address the issue described above. The changes have been analyzed and implemented to resolve the reported problem.

### Review Notes
- 🤖 This PR was created by Claude Code Background Agent
- 📋 Please review the changes carefully before merging
- 🧪 Test the changes in your development environment
- 🔍 Verify that the original issue is resolved

Closes #123
```

## Configuration

### Environment Variables

The same environment variables from Milestone 3 apply, with git configuration options:

```bash
# Required (from previous milestones)
export GITHUB_TOKEN="your_github_token"
export ANTHROPIC_API_KEY="your_anthropic_api_key"
export REPO_OWNER="username"
export REPO_NAME="repository"

# Optional git configuration
export CLAUDECLI_GIT_NAME="Claude Code Bot"
export CLAUDECLI_GIT_EMAIL="claude-code-bot@example.com"
```

### Command-Line Options

```bash
# PR creation control
--create-pr=true/false    # Enable/disable PR creation (default: true)

# Git configuration
--git-name="Name"         # Git user name for commits
--git-email="email"       # Git user email for commits
```

## Usage Examples

### Basic Issue Processing with PR Creation

```bash
# Process issue and create PR
claudecli issue 42

# Output:
# Processing issue #42 for user/repo
# Fetching issue #42 from user/repo...
# Processing issue: Fix broken navigation menu
# Cloning repository to cache...
# Creating Docker container...
# Starting container...
# Executing Claude Code...
# Claude Code execution completed successfully
# Checking for repository changes...
# Changes detected, proceeding with PR creation...
# Creating branch: claude-code/fix-issue-42
# Committing changes...
# Pushing branch to GitHub...
# Creating pull request...
# ✅ Successfully created pull request: https://github.com/user/repo/pull/43
# 📋 PR #43: Fix #42: Fix broken navigation menu
```

### Skip PR Creation

```bash
# Only run Claude Code without creating PR
claudecli issue 42 --create-pr=false
```

### Custom Git Configuration

```bash
# Use custom git credentials
claudecli issue 42 --git-name="Development Bot" --git-email="dev@company.com"
```

## Error Handling

The system handles various error conditions:

- **No Changes**: If Claude Code doesn't make changes, PR creation is skipped
- **Git Errors**: Detailed error messages for git operations
- **GitHub API Errors**: Comprehensive error handling for PR creation failures
- **Branch Conflicts**: Automatic branch name generation to avoid conflicts

## Integration with Previous Milestones

Milestone 4 builds on all previous milestones:

- **M1**: CLI framework and configuration
- **M2**: Docker container management
- **M3**: Claude Code execution and logging
- **M4**: Complete GitHub integration with PR workflow

## Testing

Test the complete workflow:

```bash
# Test with a real issue
claudecli issue 123

# Verify PR creation in GitHub
# Check that labels are applied correctly
# Confirm issue is linked to PR
```

## Next Steps: Milestone 5

The foundation is now ready for Milestone 5: "Multi-issue processing & batching", which will add capabilities to:

- Process multiple issues in batches
- Handle dependencies between issues
- Optimize container reuse for efficiency
- Provide batch status reporting

## Technical Details

### Branch Naming Convention

Branches are created with the format: `claude-code/fix-issue-{number}`

### Commit Message Format

```
Fix #{number}: {issue_title}

Automatically generated fix by Claude Code
```

### Error Recovery

- Failed PR creation attempts provide detailed logs
- Git operations are atomic where possible
- Container cleanup happens regardless of PR creation success

## Architecture Impact

Milestone 4 completes the core issue-to-PR workflow, providing:

- ✅ Full automation from issue to pull request
- ✅ Rich metadata and linking
- ✅ Intelligent categorization via labels
- ✅ Professional PR presentation
- ✅ Robust error handling and recovery

The system now provides a complete, production-ready workflow for automated issue resolution with proper GitHub integration.
