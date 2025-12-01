#!/usr/bin/env bash
set -euo pipefail

# Bootstrap goneat and foundation tools (yamlfmt, prettier, etc.)
#
# Pattern (v0.3.9+):
#   1. Download goneat binary with SHA256 verification
#   2. Initialize goneat tools config if needed (goneat doctor tools init)
#   3. Install foundation tools (goneat doctor tools --scope foundation --install)
#      - goneat auto-installs package managers (bun/brew) if needed
#      - Then installs tools via the package manager
#
# To update: change GONEAT_VERSION and corresponding SHA256 checksums

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_DIR="${ROOT_DIR}/bin"

GONEAT_VERSION="v0.3.9"

# Map uname outputs to archive naming
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH_RAW="$(uname -m)"
case "${ARCH_RAW}" in
  x86_64) ARCH="amd64" ;;
  aarch64 | arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: ${ARCH_RAW}" >&2; exit 1 ;;
esac

ARCHIVE="goneat_${GONEAT_VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/fulmenhq/goneat/releases/download/${GONEAT_VERSION}/${ARCHIVE}"
case "${OS}-${ARCH}" in
  darwin-amd64)
    EXPECTED_SHA="3a054db2d58d5a4f7a3f7fb9f8d5fba4a92c9495e9ba03ced136fcbc91be7866"
    ;;
  darwin-arm64)
    EXPECTED_SHA="830850afe860ec3773f5cc9f9eb693e3bb6aa6b8fd5bd30bcd54516d843d3d5a"
    ;;
  linux-amd64)
    EXPECTED_SHA="2541a8d75c565ff4cca71fd090110e7ae0acaa919e5ff8c2cbcd382110a67618"
    ;;
  linux-arm64)
    EXPECTED_SHA="15f49a33958c114916d9c7965ef3ace0b971855f626eb110abe58a9a7eae1d1b"
    ;;
  *)
    EXPECTED_SHA=""
    ;;
esac

mkdir -p "${BIN_DIR}"

# Download goneat if not present or if version mismatch
NEED_DOWNLOAD=false
if [[ ! -x "${BIN_DIR}/goneat" ]]; then
  NEED_DOWNLOAD=true
elif ! "${BIN_DIR}/goneat" version 2>/dev/null | grep -q "${GONEAT_VERSION}"; then
  echo "Upgrading goneat to ${GONEAT_VERSION}..."
  NEED_DOWNLOAD=true
fi

if [[ "${NEED_DOWNLOAD}" == "true" ]]; then
  echo "Downloading goneat ${GONEAT_VERSION} for ${OS}/${ARCH}..."
  tmpfile="$(mktemp)"
  trap 'rm -f "${tmpfile}"' EXIT

  curl -sSL -o "${tmpfile}" "${URL}"

  # Verify checksum if available (skip if PENDING - release not yet published)
  if [[ -n "${EXPECTED_SHA}" && "${EXPECTED_SHA}" != "PENDING" ]]; then
    ACTUAL_SHA="$(shasum -a 256 "${tmpfile}" | cut -d' ' -f1)"
    if [[ "${ACTUAL_SHA}" != "${EXPECTED_SHA}" ]]; then
      echo "❌ Checksum mismatch!" >&2
      echo "   Expected: ${EXPECTED_SHA}" >&2
      echo "   Actual:   ${ACTUAL_SHA}" >&2
      exit 1
    fi
    echo "✅ Checksum verified"
  fi

  tar -xzf "${tmpfile}" -C "${BIN_DIR}"
  chmod +x "${BIN_DIR}/goneat"
  echo "goneat installed to ${BIN_DIR}/goneat"
else
  echo "goneat ${GONEAT_VERSION} already present at ${BIN_DIR}/goneat"
fi

# Initialize goneat tools config if not present (v0.3.7+ requirement)
# This creates .goneat/tools.yaml in goneat's standard format
if [[ ! -f "${ROOT_DIR}/.goneat/tools.yaml" ]] || ! grep -q "^scopes:" "${ROOT_DIR}/.goneat/tools.yaml" 2>/dev/null; then
  echo "Initializing goneat doctor tools config..."
  "${BIN_DIR}/goneat" doctor tools init --force
fi

# Install foundation tools (yamlfmt, prettier, etc.) via goneat doctor
# v0.3.9+: goneat auto-installs bun/brew if needed, then installs tools
echo "Installing foundation tools via goneat doctor..."
"${BIN_DIR}/goneat" doctor tools --scope foundation --install --yes

echo "✅ Bootstrap complete"
