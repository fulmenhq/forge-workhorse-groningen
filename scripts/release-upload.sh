#!/usr/bin/env bash

set -euo pipefail

TAG="${1:-}"
SOURCE_DIR="${2:-dist/release}"

if [[ -z "${TAG}" ]]; then
    echo "usage: $0 vX.Y.Z [source_dir]" >&2
    exit 1
fi

if ! command -v gh > /dev/null 2>&1; then
    echo "❌ gh (GitHub CLI) not found in PATH" >&2
    echo "Install: https://cli.github.com/" >&2
    exit 1
fi

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
KEY_FILE="${ROOT_DIR}/dist/signing/public-key.asc"

if [[ ! -d "${SOURCE_DIR}" ]]; then
    echo "❌ Source dir not found: ${SOURCE_DIR}" >&2
    exit 1
fi

assets=("${SOURCE_DIR}"/*.asc)
if [[ ${#assets[@]} -eq 0 ]]; then
    echo "❌ No signature assets found in ${SOURCE_DIR} (expected *.asc)" >&2
    echo "Run signing first: SIGNING_KEY_ID=... scripts/sign-release-artifacts.sh ${SOURCE_DIR}" >&2
    exit 1
fi

if [[ ! -f "${KEY_FILE}" ]]; then
    echo "❌ Public key not found: ${KEY_FILE}" >&2
    echo "Expected signing step to export it." >&2
    exit 1
fi

echo "→ Uploading signatures and public key to ${TAG} (clobber)"
gh release upload "${TAG}" "${SOURCE_DIR}"/*.asc "${KEY_FILE}" --clobber

echo "✅ Upload complete"
