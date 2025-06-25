# Setup Script Guide

## Overview

Setup scripts allow you to customize the Claude Code development environment by installing additional tools, configuring settings, and preparing the container for specific project requirements.

## Configuration Methods

### 1. Environment Variable (Global)

```bash
export SETUP_SCRIPT=/path/to/your/setup.sh
claudecli issue run --number 123
```

### 2. Command Line Flag (Per-run)

```bash
claudecli issue run --number 123 --setup-script ./project-setup.sh
```

*Note: Command line flag takes precedence over environment variable.*

## Sample Setup Scripts

### Basic Development Tools

```bash
#!/bin/bash
set -e

echo "Installing essential development tools..."

# Update package manager
apt-get update -qq

# Install search and navigation tools
apt-get install -y ripgrep fd-find fzf

# Install build tools
apt-get install -y build-essential cmake ninja-build

# Clean up
apt-get clean
rm -rf /var/lib/apt/lists/*
```

### Node.js Project Setup

```bash
#!/bin/bash
set -e

echo "Setting up Node.js development environment..."

# Install global Node.js tools
npm install -g typescript ts-node @types/node
npm install -g eslint prettier jest

# Install project dependencies if package.json exists
if [ -f "/workspace/package.json" ]; then
    cd /workspace
    npm install
fi
```

### Python Project Setup

```bash
#!/bin/bash
set -e

echo "Setting up Python development environment..."

# Install Python development tools
pip3 install --upgrade pip
pip3 install black isort flake8 mypy pytest

# Install project dependencies if requirements.txt exists
if [ -f "/workspace/requirements.txt" ]; then
    cd /workspace
    pip3 install -r requirements.txt
fi
```

### Rust Project Setup

```bash
#!/bin/bash
set -e

echo "Setting up Rust development environment..."

# Install Rust tools
cargo install --locked ast-grep
cargo install ripgrep fd-find

# Install Rust components
rustup component add clippy rustfmt

# Build project if Cargo.toml exists
if [ -f "/workspace/Cargo.toml" ]; then
    cd /workspace
    cargo fetch
fi
```

### Multi-language Setup

```bash
#!/bin/bash
set -e

echo "Setting up multi-language development environment..."

# Update system
apt-get update -qq

# Essential tools
apt-get install -y ripgrep fd-find fzf tree jq curl wget

# Build tools
apt-get install -y build-essential cmake ninja-build

# Node.js tools
npm install -g typescript eslint prettier

# Python tools
pip3 install black isort flake8 mypy

# Rust tools
cargo install ast-grep

# Git tools
apt-get install -y git-delta tig

# Cleanup
apt-get clean
rm -rf /var/lib/apt/lists/*

echo "Development environment ready!"
```

## Best Practices

### Script Structure

```bash
#!/bin/bash
set -e  # Exit on any error

echo "🔧 Starting setup..."

# Your setup commands here

echo "✅ Setup complete!"
```

### Error Handling

```bash
#!/bin/bash
set -e

# Function for error handling
handle_error() {
    echo "❌ Setup failed at line $1"
    exit 1
}
trap 'handle_error $LINENO' ERR

# Your setup commands here
```

### Conditional Installation

```bash
#!/bin/bash
set -e

# Install tools only if not already present
if ! command -v rg &> /dev/null; then
    echo "Installing ripgrep..."
    apt-get update -qq
    apt-get install -y ripgrep
fi

# Project-specific setup
if [ -f "/workspace/package.json" ]; then
    echo "Node.js project detected, installing dependencies..."
    cd /workspace && npm install
fi
```

### Package Manager Updates

```bash
#!/bin/bash
set -e

# Always update package lists for apt
echo "Updating package lists..."
apt-get update -qq

# Install packages
apt-get install -y your-packages

# Clean up to reduce image size
apt-get clean
rm -rf /var/lib/apt/lists/*
```

## Available Tools in Base Image

The Claude Code base image includes:

- Node.js 20+
- Python 3.8+
- Rust (latest stable)
- Git
- curl, wget
- Basic build tools
- Claude CLI

## Environment Variables

Your setup script has access to issue-related environment variables:

- `ISSUE_NUMBER`: GitHub issue number
- `ISSUE_TITLE`: Issue title
- `REPO_OWNER`: Repository owner
- `REPO_NAME`: Repository name

Example usage:

```bash
#!/bin/bash
echo "Setting up environment for issue #$ISSUE_NUMBER: $ISSUE_TITLE"
echo "Repository: $REPO_OWNER/$REPO_NAME"

# Conditional setup based on repository
if [ "$REPO_NAME" == "my-nodejs-project" ]; then
    npm install -g typescript
fi
```

## Working Directory

Setup scripts run with `/workspace` as the working directory, where your repository is mounted.

## Troubleshooting

### Script Not Found

- Ensure the script path is absolute or relative to current working directory
- Verify the script file exists and is readable

### Permission Denied

- Make sure the script is executable: `chmod +x setup.sh`
- The CLI will automatically set execute permissions inside the container

### Package Installation Failures

- Always run `apt-get update` before installing packages
- Use `-y` flag for non-interactive installation
- Consider using `apt-get update -qq` to reduce output

### Network Issues

- Some corporate networks may block package downloads
- Consider using cached packages or private registries

## Integration with CI/CD

Setup scripts work well in CI environments:

```yaml
# GitHub Actions example
- name: Run Claude Code Agent
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    SETUP_SCRIPT: ./.github/scripts/claude-setup.sh
  run: |
    claudecli issue run --number ${{ github.event.issue.number }}
```

This allows you to maintain project-specific development environments while leveraging Claude Code's capabilities.
