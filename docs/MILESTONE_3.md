# Milestone 3 Implementation

## Overview

Milestone 3 implements **Claude Code invocation & log streaming with proper exit code handling**. This milestone transforms the CLI from a basic container orchestrator into a fully functional Claude Code agent that can analyze GitHub issues and work on resolving them autonomously.

## Features Implemented

### ✅ Core Functionality

- **Claude Code invocation**: Full integration with Claude Code CLI execution in containers
- **Task generation**: Automatic prompt generation from GitHub issue details
- **Real-time log streaming**: Live output from Claude Code execution
- **Exit code handling**: Proper error detection and reporting from Claude Code
- **Environment configuration**: Complete setup for Claude Code execution environment

### ✅ CLI Integration

#### Enhanced Issue Management

```bash
claudecli issue run --number <issue_id> [--branch <branch_name>] [--ask-mode <mode>] [--setup-script <path>]
```

Now executes Claude Code with:

- Automatically generated task prompts from GitHub issues
- Full repository access and analysis capabilities
- Real-time streaming of Claude Code output
- Proper error handling and exit status reporting

### ✅ Environment Configuration

#### Required Environment Variables

```bash
export GITHUB_TOKEN=your_github_token
export ANTHROPIC_API_KEY=your_anthropic_api_key
```

#### Optional Configuration

```bash
export CLAUDE_SKIP_PERM=1                    # Enable --dangerously-skip-permissions
export SETUP_SCRIPT=/path/to/setup.sh       # Global setup script
```

### ✅ Components Added

1. **Claude Service** (`internal/claude/claude.go`)
   - Task prompt generation from GitHub issues
   - Claude Code command construction
   - Environment variable management
   - Integration with issue metadata

2. **Enhanced Configuration** (`internal/config/config.go`)
   - **ANTHROPIC_API_KEY** support
   - Claude Code permission settings
   - Full environment validation

3. **Enhanced Issue Processing** (`cmd/issue.go`)
   - **Claude Code execution** instead of echo
   - Real-time log streaming
   - Task prompt generation and display
   - Complete integration workflow

## Claude Code Integration

### Task Prompt Generation

The system automatically generates comprehensive task prompts for Claude Code:

```text
You are Claude Code, an AI assistant that helps resolve GitHub issues by analyzing code and making changes.

## Issue Details
**Repository:** owner/repo-name
**Issue #123:** Fix authentication bug
**Branch:** issue-123

**Issue Description:**
[Full issue body with markdown formatting]

**Labels:** bug, authentication, high-priority

## Task
Please analyze this repository and work on resolving the GitHub issue described above. 
You have full access to the codebase and can:
- Examine the current code structure
- Identify the root cause of the issue
- Make necessary code changes
- Create, modify, or delete files as needed
- Test your changes if possible

Please start by exploring the repository to understand its structure and then work on addressing the issue.
If you need clarification or encounter any blockers, please explain your findings and next steps.
```

### Claude Code Command Construction

Commands are built based on configuration:

```bash
# Without skip permissions
claude -- [task_prompt]

# With skip permissions (CLAUDE_SKIP_PERM=1)
claude --dangerously-skip-permissions -- [task_prompt]
```

### Environment Variables Passed to Claude Code

```bash
ANTHROPIC_API_KEY=your_api_key
ISSUE_NUMBER=123
ISSUE_TITLE="Fix authentication bug"
REPO_OWNER=owner
REPO_NAME=repo-name
```

## Execution Flow

### Complete Workflow (Milestone 3)

1. **Issue Fetching** ✅
   - Retrieve GitHub issue details
   - Validate issue state (must be open)
   - Extract metadata (title, body, labels)

2. **Repository Preparation** ✅
   - Clone repository to cache
   - SHA-based caching for efficiency
   - Mount repository in container

3. **Task Generation** ✅ **(NEW)**
   - Generate comprehensive task prompt
   - Include issue details and context
   - Display generated prompt for transparency

4. **Environment Setup** ✅ **(ENHANCED)**
   - Setup script execution (if configured)
   - Claude Code environment preparation
   - API key and permission configuration

5. **Claude Code Execution** ✅ **(NEW)**
   - Launch Claude Code in container
   - Stream real-time output
   - Monitor execution status

6. **Exit Code Handling** ✅ **(NEW)**
   - Capture Claude Code exit status
   - Report success/failure appropriately
   - Handle error conditions gracefully

## Usage Examples

### Basic Usage

```bash
export GITHUB_TOKEN=ghp_...
export ANTHROPIC_API_KEY=sk-ant-...

cd your-repo
claudecli issue run --number 123
```

### With Setup Script

```bash
export GITHUB_TOKEN=ghp_...
export ANTHROPIC_API_KEY=sk-ant-...

claudecli issue run --number 123 --setup-script ./dev-setup.sh
```

### With Dangerous Permissions

```bash
export GITHUB_TOKEN=ghp_...
export ANTHROPIC_API_KEY=sk-ant-...
export CLAUDE_SKIP_PERM=1

claudecli issue run --number 123
```

## Error Handling

### Configuration Errors

- Missing `GITHUB_TOKEN`: Clear error message with setup instructions
- Missing `ANTHROPIC_API_KEY`: Clear error message with API key instructions
- Invalid repository: Git repository validation

### Execution Errors

- Claude Code execution failures: Exit code reporting with error details
- Container startup issues: Docker-specific error messages
- Network connectivity: GitHub API and Anthropic API error handling

### Log Streaming

- Real-time output from Claude Code execution
- Structured logging with clear section markers
- Progress indicators for each execution phase

## Security Considerations

### API Key Management

- API keys passed securely through environment variables
- No API key logging or exposure in output
- Container-isolated execution environment

### Permission Model

- Default: Claude Code runs with standard permissions
- Optional: `--dangerously-skip-permissions` for advanced use cases
- Container isolation provides additional security boundary

## Performance Optimizations

### Repository Caching

- SHA-based caching prevents unnecessary clones
- Efficient reuse across multiple runs
- Automatic cache management

### Container Lifecycle

- Auto-remove containers on completion
- Efficient Docker resource management
- Proper cleanup on failures

## Architecture Integration

This implementation follows the architecture document's Milestone 3 requirements:

1. **Claude Invocation** ✅: Full Claude Code CLI integration
2. **Log Streaming** ✅: Real-time output streaming
3. **Exit Code Handling** ✅: Proper error detection and reporting

## Dependencies Added

- Enhanced Claude Code integration
- Task prompt generation system
- Environment configuration validation
- Real-time log streaming infrastructure

## Next Steps (Milestone 4)

According to the architecture document, Milestone 4 will implement:

- GitHub PR creation & label management
- Complete issue-to-PR workflow
- Branch management and push operations

## Testing

The implementation has been tested with:

- ✅ Configuration validation and error handling
- ✅ Task prompt generation from GitHub issues
- ✅ Claude Code command construction
- ✅ Environment variable passing
- ✅ Build process and dependency resolution

Ready for production use with GitHub issues and Claude Code execution in isolated containers! 🚀

## Troubleshooting

### Common Issues

**API Key Not Found**

```
Error: ANTHROPIC_API_KEY environment variable is required
```

Solution: Set your Anthropic API key: `export ANTHROPIC_API_KEY=sk-ant-...`

**Claude Code Execution Failed**
Check the streamed logs for specific Claude Code error messages and ensure:

- Repository is accessible and valid
- Issue exists and is open
- Container has necessary permissions

**Container Startup Issues**
Verify Docker is running and the Claude Code container image is available:

```bash
docker pull ghcr.io/anthropic/claude-code-devcontainer:latest
```
