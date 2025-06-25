# Milestone 2 Implementation

## Overview

Milestone 2 implements the minimal `issue run` functionality with repository cloning, Docker container management, echo capabilities, and **setup script support** to ensure Claude Code has a properly configured development environment.

## Features Implemented

### ✅ Core Functionality

- **Repository cloning**: Shallow clone to local cache (`$XDG_CACHE_HOME/claudecli/repos/<owner>-<name>/<sha>`)
- **Docker container management**: Start, stop, list, and remove containers
- **Container execution**: Run containers with mounted repository and environment variables
- **Log streaming**: Real-time container output streaming
- **Echo demonstration**: Simple command execution to verify container functionality
- **Setup script support**: Configurable script execution for custom development environment setup

### ✅ CLI Commands

#### Issue Management

```bash
claudecli issue run --number <issue_id> [--branch <branch_name>] [--ask-mode <mode>] [--setup-script <path>]
```

#### Container Management

```bash
claudecli container ls                    # List all containers
claudecli container stop <container_id>   # Stop and remove a container
```

### ✅ Setup Script Configuration

The CLI now supports setup scripts to ensure Claude Code has access to necessary development tools:

#### Environment Variable

```bash
export SETUP_SCRIPT=/path/to/your/setup.sh
```

#### Command Line Flag

```bash
claudecli issue run --number 123 --setup-script ./setup.sh
```

#### Sample Setup Script (`setup.sh`)

```bash
#!/bin/bash
# Install development tools for Claude Code environment
apt-get update -qq
apt-get install -y ripgrep fd-find fzf build-essential cmake ninja-build
cargo install --locked ast-grep
npm install -g typescript ts-node @types/node
pip3 install black isort flake8 mypy pytest
```

### ✅ Components Added

1. **Docker Service** (`internal/docker/docker.go`)
   - Container lifecycle management
   - Image pulling with progress output
   - Volume mounting for repository access
   - **Setup script mounting and execution**
   - Log streaming and status monitoring

2. **Enhanced Git Operations** (`internal/git/git.go`)
   - Cache directory management
   - Repository cloning with reuse logic
   - SHA-based caching for efficiency

3. **Container Commands** (`cmd/container.go`)
   - List running containers with formatted output
   - Stop and remove containers safely

4. **Enhanced Configuration** (`internal/config/config.go`)
   - **Setup script path configuration**
   - Environment variable support

## Architecture

The implementation follows the architecture document's Milestone 2 requirements with enhanced development environment support:

1. **Repository Handling**: Local shallow clone to cache before container execution
2. **Docker Integration**: Full container lifecycle management using Docker SDK
3. **Isolation**: Containers run with mounted repository and isolated environment
4. **Logging**: Real-time log streaming from container execution
5. **Development Environment**: Configurable setup scripts for tool installation

## Development Environment Setup

The setup script system addresses the need for Claude Code to operate in a properly supported environment:

### Base Layer (Docker Image)

Expected to include foundational tools:

- Node.js 20+
- Python 3.8+
- Git
- Build essentials
- Claude CLI

### Custom Layer (Setup Script)

Users can specify additional tools:

- `fd`, `rg`, `fzf` for fast searching
- `ast-grep` for AST-based code analysis
- `cmake`, `ninja` for build systems
- Language-specific tools and linters
- Project-specific dependencies

## Usage Example

```bash
# Set required environment variable
export GITHUB_TOKEN=your_github_token

# Optional: Set global setup script
export SETUP_SCRIPT=./setup.sh

# Navigate to a git repository
cd /path/to/your/repo

# Run with custom setup script
./claudecli issue run --number 123 --setup-script ./project-setup.sh

# List containers
./claudecli container ls

# Stop a container if needed
./claudecli container stop <container_id>
```

## Dependencies Added

- `github.com/docker/docker/client` - Docker SDK client
- `github.com/docker/docker/api/types/*` - Docker API types
- Enhanced error handling and logging
- Setup script execution support

## Next Steps (Milestone 3)

According to the architecture document, Milestone 3 will implement:

- Claude invocation & log streaming
- Handle exit codes
- Full integration with Claude Code agent execution
- **Enhanced development environment with pre-configured tools**

## Testing

The implementation has been tested with:

- ✅ CLI command structure and help output
- ✅ Docker service connectivity
- ✅ Container listing functionality
- ✅ Build process and dependency resolution
- ✅ Setup script configuration and execution
- ✅ Environment variable and flag support

Ready for integration testing with real GitHub issues, Docker containers, and custom development environments.
