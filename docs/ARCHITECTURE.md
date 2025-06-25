# Claude Code Agent CLI – Architecture

## 1. System Overview

The CLI coordinates three primary actors:

1. **Host process** – the Go binary executed by the user or CI.
2. **Docker daemon** – launches pre-hardened dev-containers.
3. **Claude Code model** – runs inside the container via the `claude` CLI, using the repository mounted within.

## 2. Components

| Component      | Responsibility |
|----------------|----------------|
| **CLI binary** | Parses commands, loads config, calls GitHub & Docker APIs, streams logs |
| **Dev-container image** | Provides Node.js 20, Claude CLI, hardened firewall |
| **Supervisor script** | Started inside container; invokes Claude, monitors exit status |
| **GitHub service** | Issue & PR operations via `go-github` |

---

### Repository Handling Strategy

Prior to **clarification**, the CLI now performs a **local shallow clone** of the
repository (or re-uses an existing checkout in a cache directory). This gives
the host-side Claude process full read-only visibility into the codebase while
questions are formulated.  The clone happens in `$XDG_CACHE_HOME/claudecli/repos/<owner>-<name>/<sha>`
and can be reused across runs.

During the **build phase** the container still performs its own fresh clone to
maintain isolation.  Mounting the host checkout is possible with
`--mount-mode host` but, for security, the default remains `clone`.

`git worktree` **remains available** via `--mount-mode worktree` for extremely
large monorepos where cloning would be prohibitive.  The trade-offs are:

| Aspect            | Clone in container (default) | Host worktree mount |
|-------------------|------------------------------|---------------------|
| Isolation         | Full                         | Leaks host FS paths |
| Setup complexity  | Low                          | Requires cleanup    |
| Initial speed     | Slightly slower (network)    | Fast                |
| Large repos       | Can use sparse checkout      | Generally faster    |
| Security risk     | Minimal                      | Higher              |

---

### Clarification Gate

The CLI **never launches or mounts a container until the clarifying-questions
phase is resolved**:

1. The host process fetches the issue and runs a small LLM prompt (cost-effective
   model) to decide if more information is required.
2. If questions are needed, it posts a comment and exits with a non-zero status
   (`exit 2`). No Docker activity occurs.
3. Once the issue is updated with answers, the operator/CI reruns `issue run`.
   Only then does the CLI start the container and perform the clone/mount.

This gate avoids burning compute on containers that may immediately stall
waiting for human input.

---

## 3. CLI Surface

```bash
claudecli [global-flags] command [command-flags]

Global flags
  --log-level      debug|info|warn|error

Commands
  issue run        Start a container for the current Git repo's issue
     --number      issue number
     --branch      working branch name (default: issue-<num>)
     --ask-mode    auto|interactive|none
  container ls     List running agent containers
  container stop   Stop & remove a container
  pr open          Open a PR for an existing branch
  auth login       Store GitHub token locally (optional)
```

---

### Environment Variables

Configuration is **environment-first**. No YAML files are required.

| Variable          | Purpose                                       |
|-------------------|-----------------------------------------------|
| `GITHUB_TOKEN`    | Personal access token with `repo` scope        |
| `CLAUDE_SKIP_PERM`| If set, passes `--dangerously-skip-permissions`|
| `CACHE_DIR`       | Overrides default `$XDG_CACHE_HOME/claudecli`  |

If `GITHUB_TOKEN` is missing the CLI exits with `status 1` and instructs the
user to export it.

---

## 4. Detailed Workflow

### End-to-End Flow with Clarification Stage (revised)

1. **Fetch issue** – CLI retrieves metadata.
2. **Local repo clone** – CLI performs a shallow clone of the target repo and
   checks out the default branch.
3. **Local static exploration** – Host-side Claude (safe mode) analyses the code
   and `CLAUDE.md` to build context.
4. **Question generation** – Claude posts clarifying questions on the issue.
5. **Turn-by-turn resolution** – Loop until Claude signals all questions
   answered.
6. **Task specification** – Claude writes **`TASK.md`** to the cloned repo and
   the CLI commits it to a new branch (`issue-<num>-plan`).
7. **Dev-container execution** – CLI launches container, shallow-clones the same
   branch, and starts Claude in unrestricted mode with `TASK.md` as guidance.
8. **Pull request** – On success, branch is pushed & PR opened.
9. **Container teardown** – Logs saved, container removed.

This two-phase design (clarify → build) minimises wasted container minutes and provides an auditable plan for every automated change.

### Clarification Gate (updated)

* The container **only** launches after `TASK.md` exists **and** is committed to
  the branch.
* Host-side analysis uses the local checkout; no Docker usage until build
  phase.

---

## 5. Configuration File (`~/.claudecli.yaml`)

```yaml
github_token: ghp_xxx
anthropic_api_key: sk-ant-xxx
docker:
  image: ghcr.io/anthropic/claude-code-devcontainer:latest
  mount_mode: worktree   # worktree | clone
  remove_on_exit: true
agent:
  max_tokens: 8192
  model: claude-3-code
  ask_mode: auto         # see --ask-mode override
repositories:
  - owner: myorg
    name: big-monorepo
    default_branch: main
```

---

## 6. Security & Isolation

* Uses Anthropic's dev-container firewall (default-deny outbound policy).
* Host secrets are injected only for the container's lifetime.
* `--dangerously-skip-permissions` is safe within the container boundary.

---

## 7. Implementation Roadmap

| Milestone | Deliverable | Notes |
|-----------|-------------|-------|
| M1 | Project scaffolding, config loader, `auth login` | `cobra` CLI in Go |
| M2 | Minimal `issue run`: clone repo, start container, echo | Docker SDK |
| M3 | Claude invocation & log streaming | Handle exit codes |
| M4 | GitHub PR creation & label management | go-github |
| M5 | Clarification decision model | Pluggable |
| M6 | Worktree vs clone optimisation, detaching, retries | Polish |
| M7 | Tests, CI, docs | Release v1.0 |

---

## 8. Future Enhancements

1. Kubernetes executor for parallel issues.
2. Slack / Discord notifications.
3. Dependency caching layer.
4. Additional policy engine (seccomp, AppArmor).

---

## 9. Open Questions

1. Strategy for clarification detection (regex vs. LLM judgement)?
2. Should long-running containers self-update on new pushes?
3. How to cancel an in-progress agent if the issue is closed?
