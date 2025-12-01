#!/usr/bin/env bash
set -euo pipefail

# Bootstrap goneat into ./bin from the pinned version in .goneat/tools.yaml
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_DIR="${ROOT_DIR}/bin"
MANIFEST="${ROOT_DIR}/.goneat/tools.yaml"

GONEAT_VERSION="v0.3.8"

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

mkdir -p "${BIN_DIR}"

if [[ -x "${BIN_DIR}/goneat" ]]; then
  echo "goneat already present at ${BIN_DIR}/goneat"
  exit 0
fi

echo "Downloading goneat ${GONEAT_VERSION} for ${OS}/${ARCH}..."
tmpfile="$(mktemp)"
trap 'rm -f "${tmpfile}"' EXIT

curl -sSL -o "${tmpfile}" "${URL}"
tar -xzf "${tmpfile}" -C "${BIN_DIR}"
chmod +x "${BIN_DIR}/goneat"

echo "goneat installed to ${BIN_DIR}/goneat"
