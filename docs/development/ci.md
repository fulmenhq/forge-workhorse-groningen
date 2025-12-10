# CI/CD Configuration

This document explains the CI/CD setup for this repository.

## Container-Based CI Pattern

This repository uses the **goneat-tools container** (`ghcr.io/fulmenhq/goneat-tools:latest`) for CI jobs. This is the recommended "low friction" approach from goneat v0.3.14+.

### Why Containers?

The container provides all foundation tools pre-installed:

- `prettier` - Markdown/JSON formatting
- `yamlfmt` - YAML formatting
- `jq` / `yq` - JSON/YAML processing
- `rg` (ripgrep) - Fast search
- `curl` / `wget` - HTTP tools

This eliminates tool installation friction in CI - no package manager setup, no version conflicts, no install failures.

### Container Permissions (`--user root`)

**IMPORTANT**: Do not remove the `options: --user root` from container configurations.

```yaml
container:
  image: ghcr.io/fulmenhq/goneat-tools:latest
  # REQUIRED: actions/checkout@v4 needs write access to /__w/_temp/_runner_file_commands/
  options: --user root
```

#### Why is this required?

GitHub Actions' `actions/checkout@v4` writes state files to the runner's temp directory at `/__w/_temp/_runner_file_commands/`. When running inside a container, the default user may not have write permissions to this directory.

Without `--user root`, you'll see errors like:

```
Error: EACCES: permission denied, open '/__w/_temp/_runner_file_commands/save_state_...'
```

#### Is this a security concern?

No. The container is ephemeral and isolated to this CI job. Running as root inside the container does not grant elevated permissions on the GitHub runner host. This is standard practice for GitHub Actions container jobs.

### Runner Temp Permissions

Recent hardening of `ghcr.io/fulmenhq/goneat-tools` runs the image as a non-root user by default. Even with `options: --user root`, GitHub Actions initializes the `_runner_file_commands` directory before steps execute, which can leave `/__w/_temp` owned by a different uid/gid. To avoid `EACCES` failures when saving state, add a step immediately after checkout that relaxes permissions on that directory:

```yaml
- name: Fix temp permissions
  run: |
    set -euo pipefail
    sudo install -d -m 0777 /__w/_temp || true
    sudo install -d -m 0777 /__w/_temp/_runner_file_commands || true
    sudo chown -R $(id -u):$(id -g) /__w/_temp || true
    sudo chmod -R 777 /__w/_temp || true
```

This step is idempotent and safe to run in every container job. Keep it before any step that writes workflow state (e.g., `actions/cache`, `save-state`, `set-output`).

### Additional Hardening Patterns

#### Minimize `GITHUB_TOKEN` capabilities

Explicitly declare the workflow `permissions` block (for this template: `contents: read`) so the implicit token cannot mutate repository state even if a step is compromised. This keeps the example aligned with GitHub's least-privilege guidance.

#### Avoid persisting checkout credentials

Pass `persist-credentials: false` to `actions/checkout@v4`. CI jobs in this template never push, so there's no reason to store the short-lived token inside `.git/config`. Downstream users can override when they need to push tags or release artifacts.

#### Enforce strict shell options in scripts

Add `set -euo pipefail` at the top of every multi-line `run` script. This catches unset variables, stops on the first failing command, and prevents silent formatting or build failures inside the container.

### CI Jobs

1. **format-check**: Validates formatting using container tools (yamlfmt, prettier)
2. **build-test**: Builds and tests the application using container tools + goneat binary

### Local Development

For local development, you have two options:

1. **Use the container** (recommended for consistency):
   ```bash
   docker run --rm -v "$(pwd)":/work -w /work --entrypoint "" \
     ghcr.io/fulmenhq/goneat-tools:latest yamlfmt -lint .
   ```

2. **Install tools via goneat**:
   ```bash
   ./scripts/install-goneat.sh
   ./bin/goneat doctor tools --scope foundation --install
   ```

## References

- [goneat-tools container](https://github.com/fulmenhq/goneat-tools)
- [goneat documentation](https://github.com/fulmenhq/goneat)
- [GitHub Actions container jobs](https://docs.github.com/en/actions/using-jobs/running-jobs-in-a-container)
