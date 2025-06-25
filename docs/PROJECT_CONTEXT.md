# Claude Code Agent CLI – Project Context

## Overview

`claude-code-agent-cli` is a command-line tool that automates the full workflow from a GitHub issue to a finished pull request using Anthropic's Claude Code model. It orchestrates disposable Docker containers running the official Claude Code dev-container image, allowing the agent to work in an isolated, permission-free environment while keeping the host clean.

The main goals are:

1. **Automated delivery** – given an open GitHub issue, the tool can ask clarifying questions, implement the change, and open a pull request.
2. **Ephemeral execution** – all heavy computation happens inside short-lived containers built from Anthropic's hardened dev-container image.
3. **Sensible defaults, scriptable interface** – suitable for CI pipelines or long-running background workers.

---

## Objectives

1. Launch disposable Docker containers with unrestricted permissions for Claude Code.
2. Mount or clone the target repository inside the container.
3. Run Claude Code to implement the requested changes.
4. Push the resulting branch to GitHub and open a pull request.
5. Optionally post clarifying questions back to the issue before work begins.

---

## Key Features

* **Container isolation** – all agent activity occurs inside a secure dev-container.
* **GitHub integration** – reads issues, posts comments, creates branches, and opens PRs.
* **Flexible repository mounting** – choose between `git worktree` for speed or a fresh clone.
* **Config-driven** – YAML configuration with overrides via CLI flags.
* **Cross-platform binary** – single static Go executable.

---

## Technology Stack

| Concern                 | Choice | Rationale |
|-------------------------|--------|-----------|
| CLI implementation      | **Go** | Fast native binary, great concurrency, simple cross-compile |
| Docker control          | `github.com/docker/docker` client | Avoids shelling out |
| GitHub API              | `github.com/google/go-github` | Mature, well-maintained |
| LLM invocation          | `claude` CLI inside container | Matches Anthropic docs |
| Config / secrets        | YAML + env vars | Human-friendly, Git-ignore-able |

---

## Scope Out of Phase 1

* Kubernetes executor (future improvement).
* Slack/Discord notifications.
* Fine-grained policy engine beyond container isolation.
