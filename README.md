# Claude Code Background Agent

A command-line tool that automates GitHub issue resolution using Claude Code in isolated Docker containers. Handles complete workflow from issue analysis to pull request creation with intelligent labeling and linking.

## 🚀 Current Status: Milestone 4 Complete

**Latest Achievement**: Full GitHub PR creation & label management with complete issue-to-PR workflow automation.

### ✅ Completed Milestones

#### Milestone 1: CLI Scaffolding & Config Loader ✅

- [x] Cobra CLI framework with auth commands
- [x] GitHub API integration and authentication
- [x] Configuration management and validation
- [x] Environment variable support

#### Milestone 2: Docker Container Management ✅

- [x] Docker SDK integration with container lifecycle management
- [x] Repository cloning with SHA-based caching
- [x] Volume mounting for repository access
- [x] Container cleanup and resource management
- [x] Setup script support for development environment customization

#### Milestone 3: Claude Code Invocation & Log Streaming ✅

- [x] Automatic task prompt generation from GitHub issues
- [x] Claude Code execution in containers with comprehensive environment setup
- [x] Real-time log streaming from Claude Code execution
- [x] Proper exit code handling and error reporting
- [x] Full integration with GitHub issue processing

#### Milestone 4: GitHub PR Creation & Label Management ✅

- [x] **Automated Pull Request Creation** - Creates PRs with rich metadata after successful Claude Code execution
- [x] **Intelligent Label Management** - Automatically applies relevant labels based on issue type (bug-fix, enhancement, documentation)
- [x] **Branch Management** - Creates properly named branches with automatic conflict resolution
- [x] **Issue Linking** - Automatically links PRs to original issues with closure references
- [x] **Professional PR Content** - Generated PRs include comprehensive descriptions, review guidelines, and automated fix explanations
- [x] **Git Workflow Integration** - Complete git operations including commit, push, and remote branch creation
- [x] **Error Recovery** - Robust error handling with detailed logging and graceful failure recovery

### 🔮 Upcoming Milestones

#### Milestone 5: Multi-issue Processing & Batching

- [ ] Batch processing multiple issues efficiently
- [ ] Container reuse optimization for improved performance
- [ ] Dependency handling between related issues
- [ ] Comprehensive batch status reporting

#### Milestone 6: LLM Router (Claude/GPT-4o/other)

- [ ] Support for multiple LLM providers
- [ ] Dynamic model selection based on issue complexity
- [ ] Cost optimization and performance routing
- [ ] Unified interface for different AI providers

#### Milestone 7: Advanced Clarification System

- [ ] Interactive clarification for ambiguous issues
- [ ] Automated follow-up question generation
- [ ] Multi-turn conversation handling for complex requirements
- [ ] Context preservation across clarification rounds

## 🛠 Installation & Setup

### Prerequisites

- **Go 1.23+** for building the CLI
- **Docker** for containerized execution
- **Git** for repository operations
- **GitHub Personal Access Token** with repo and issues permissions
- **Anthropic API Key** for Claude Code access

### Quick Start

1. **Clone and build**:

```bash
git clone https://github.com/your-username/claude-code-background-agent
cd claude-code-background-agent
go build -o claudecli
```

2. **Configure environment**:

```bash
export GITHUB_TOKEN="your_github_token"
export ANTHROPIC_API_KEY="your_anthropic_api_key"
export REPO_OWNER="your_username"
export REPO_NAME="your_repository"
```

3. **Authenticate with GitHub**:

```bash
./claudecli auth login
```

4. **Process an issue with automatic PR creation**:

```bash
./claudecli issue 42
```

## 📋 Usage

### Core Issue Processing

```bash
# Basic issue processing with PR creation (recommended)
./claudecli issue 123

# Skip PR creation (analysis only)
./claudecli issue 123 --create-pr=false

# Custom git configuration
./claudecli issue 123 --git-name="Dev Bot" --git-email="bot@company.com"

# Use custom setup script
./claudecli issue 123 --setup-script=/path/to/custom-setup.sh
```

### Container Management

```bash
# List active containers
./claudecli container list

# Stop specific container
./claudecli container stop <container-id>

# Stop all Claude Code containers
./claudecli container stop-all
```

### Authentication

```bash
# Login with GitHub
./claudecli auth login

# Check authentication status
./claudecli auth status

# Logout
./claudecli auth logout
```

## 🏗 Architecture

The system follows a modular architecture with clear separation of concerns:

```
claudecli
├── cmd/                    # CLI command implementations
│   ├── auth.go            # GitHub authentication
│   ├── container.go       # Container management
│   ├── issue.go           # Complete issue-to-PR workflow
│   └── root.go            # CLI framework setup
├── internal/
│   ├── claude/            # Claude Code integration
│   ├── config/            # Configuration management
│   ├── docker/            # Container lifecycle
│   ├── git/               # Repository operations
│   └── github/            # PR creation & label management
└── docs/                  # Milestone documentation
```

### Key Components

- **Issue Processor**: Fetches GitHub issues and orchestrates the complete workflow
- **Docker Service**: Manages container lifecycle with volume mounting and cleanup
- **Claude Integration**: Handles task generation and Claude Code execution
- **GitHub Service**: Creates pull requests with intelligent labeling
- **Git Operations**: Manages repository cloning, branching, and remote operations

## 🔧 Configuration

### Environment Variables

```bash
# Required
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
export ANTHROPIC_API_KEY="sk-ant-xxxxxxxxxxxx"
export REPO_OWNER="your-username"
export REPO_NAME="your-repository"

# Optional
export CACHE_DIR="/tmp/claudecli-cache"        # Repository cache location
export DOCKER_IMAGE="claudecode:latest"        # Docker image for execution
export SETUP_SCRIPT="/path/to/setup.sh"       # Custom setup script
export CLAUDECLI_GIT_NAME="Claude Code Bot"   # Git commit author
export CLAUDECLI_GIT_EMAIL="bot@example.com"  # Git commit email
```

### Setup Script Support

Create custom development environments with setup scripts:

```bash
#!/bin/bash
# setup.sh - Custom development tools

# Install additional tools
apt-get update && apt-get install -y \
    fd-find \
    ripgrep \
    fzf \
    nodejs \
    npm

# Setup language-specific tools
npm install -g typescript eslint prettier
pip install black pylint mypy

# Configure git (optional)
git config --global init.defaultBranch main
```

## 🎯 Features

### Milestone 4 Highlights

#### **Automated PR Creation**

- Creates pull requests automatically after successful Claude Code execution
- Rich PR descriptions with issue context and review guidelines
- Automatic branch naming with conflict resolution
- Professional formatting with emojis and structured content

#### **Intelligent Label Management**

- Base labels: `automated-fix`, `claude-code`
- Type-based labels: `bug-fix`, `enhancement`, `documentation`
- Contextual label detection from original issue labels
- Automatic label application with error handling

#### **Complete Git Workflow**

- Repository change detection before PR creation
- Automatic branch creation with standardized naming
- Comprehensive commit messages with issue references
- Remote push with authentication handling

#### **Issue Integration**

- Automatic issue linking in PR descriptions
- Issue closure references (`Closes #123`)
- Comment posting to link PRs back to issues
- State validation (only processes open issues)

### Example Workflow Output

```bash
$ ./claudecli issue 42
Processing issue #42 for user/repo
Fetching issue #42 from user/repo...
Processing issue: Fix broken navigation menu
Cloning repository to cache...
Creating Docker container...
Starting container...
Executing Claude Code...
Claude Code execution completed successfully
Checking for repository changes...
Changes detected, proceeding with PR creation...
Creating branch: claude-code/fix-issue-42
Committing changes...
Pushing branch to GitHub...
Creating pull request...
✅ Successfully created pull request: https://github.com/user/repo/pull/43
📋 PR #43: Fix #42: Fix broken navigation menu
```

## 🧪 Testing

### Verify Complete Workflow

1. **Create a test issue** in your repository
2. **Run the processor**: `./claudecli issue <issue-number>`
3. **Check results**:
   - Verify Claude Code execution completed successfully
   - Confirm changes were detected in the repository
   - Check that a PR was created with appropriate labels
   - Verify issue linking and closure reference

### Container Management Testing

```bash
# Test container lifecycle
./claudecli container list
./claudecli container stop-all

# Verify cleanup
docker ps -a | grep claudecode
```

## 🚧 Development

### Building from Source

```bash
# Clone repository
git clone https://github.com/your-username/claude-code-background-agent
cd claude-code-background-agent

# Install dependencies
go mod download

# Build CLI
go build -o claudecli

# Run tests
go test ./...
```

### Adding Custom Commands

The CLI uses Cobra framework. Add new commands in the `cmd/` directory:

```go
// cmd/newcommand.go
var newCmd = &cobra.Command{
    Use:   "new",
    Short: "Description of new command",
    RunE:  runNewCommand,
}

func init() {
    rootCmd.AddCommand(newCmd)
}
```

## 📚 Documentation

- [**Milestone 1 Documentation**](docs/MILESTONE_1.md) - CLI scaffolding and configuration
- [**Milestone 2 Documentation**](docs/MILESTONE_2.md) - Docker integration and setup scripts
- [**Milestone 3 Documentation**](docs/MILESTONE_3.md) - Claude Code execution and logging
- [**Milestone 4 Documentation**](docs/MILESTONE_4.md) - GitHub PR creation and label management
- [**Architecture Overview**](docs/ARCHITECTURE.md) - System design and component interaction
- [**Setup Script Guide**](docs/SETUP_SCRIPT_GUIDE.md) - Custom environment configuration

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙋‍♂️ Support

- **Issues**: Report bugs and request features via [GitHub Issues](https://github.com/your-username/claude-code-background-agent/issues)
- **Discussions**: Join conversations in [GitHub Discussions](https://github.com/your-username/claude-code-background-agent/discussions)
- **Documentation**: Comprehensive guides available in the `/docs` directory

---

**Built with ❤️ for automated issue resolution using Claude Code**
