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
