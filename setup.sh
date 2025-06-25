#!/bin/bash
set -e

echo "🔧 Starting Claude Code environment setup..."

# Update package lists
echo "📦 Updating package lists..."
apt-get update -qq

# Install additional development tools
echo "🛠️  Installing additional development tools..."

# Install ripgrep for fast text search
echo "  - Installing ripgrep..."
apt-get install -y ripgrep

# Install fd for fast file finding
echo "  - Installing fd..."
apt-get install -y fd-find

# Install fzf for fuzzy finding
echo "  - Installing fzf..."
apt-get install -y fzf

# Install ast-grep for AST-based code search
echo "  - Installing ast-grep..."
cargo install --locked ast-grep

# Install additional build tools
echo "  - Installing build essentials..."
apt-get install -y build-essential cmake ninja-build

# Install Node.js tools (if not already present)
echo "  - Installing Node.js tools..."
npm install -g npm@latest
npm install -g typescript ts-node @types/node

# Install Python development tools
echo "  - Installing Python development tools..."
pip3 install --upgrade pip
pip3 install black isort flake8 mypy pytest

# Install Git tools
echo "  - Installing Git tools..."
apt-get install -y git-delta tig

# Clean up
echo "🧹 Cleaning up..."
apt-get clean
rm -rf /var/lib/apt/lists/*

echo "✅ Claude Code environment setup complete!"
echo "🚀 Environment is ready for development work."

# List installed tools for verification
echo ""
echo "📋 Installed development tools:"
echo "  - ripgrep: $(rg --version | head -1)"
echo "  - fd: $(fd --version)"
echo "  - fzf: $(fzf --version)"
echo "  - ast-grep: $(sg --version)"
echo "  - cmake: $(cmake --version | head -1)"
echo "  - ninja: $(ninja --version)"
echo "  - node: $(node --version)"
echo "  - npm: $(npm --version)"
echo "  - python3: $(python3 --version)"
echo "  - pip3: $(pip3 --version)"
echo ""
