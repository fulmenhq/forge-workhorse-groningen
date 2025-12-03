#!/usr/bin/env bash
#
# Sign release artifacts with a local GPG key.
# Usage:
#   SIGNING_KEY_ID=<key-id> scripts/sign-release-artifacts.sh [artifact_dir]
# - artifact_dir defaults to ./bin
# - Generates .asc signatures next to each artifact
# - Exports public key to dist/signing/public-key.asc
#
# This is intentionally manual to keep the template CI-light. Downstream
# consumers should replace the signing key with their own after CDRL refit.

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ARTIFACT_DIR="${1:-${ROOT_DIR}/bin}"
SIGNING_KEY_ID="${SIGNING_KEY_ID:-${GPG_KEY_ID:-}}"
EXPORT_DIR="${ROOT_DIR}/dist/signing"

if [[ -z "${SIGNING_KEY_ID}" ]]; then
  echo "❌ SIGNING_KEY_ID (or GPG_KEY_ID) is required" >&2
  exit 1
fi

if ! command -v gpg >/dev/null 2>&1; then
  echo "❌ gpg not found in PATH" >&2
  exit 1
fi

if [[ ! -d "${ARTIFACT_DIR}" ]]; then
  echo "❌ Artifact directory not found: ${ARTIFACT_DIR}" >&2
  exit 1
fi

mkdir -p "${EXPORT_DIR}"

echo "🔏 Signing artifacts in ${ARTIFACT_DIR} with key ${SIGNING_KEY_ID}..."
shopt -s nullglob
artifacts=("${ARTIFACT_DIR}"/groningen* "${ARTIFACT_DIR}"/SHA256SUMS*)
if [[ ${#artifacts[@]} -eq 0 ]]; then
  echo "⚠️  No artifacts found matching groningen* or SHA256SUMS* in ${ARTIFACT_DIR}"
fi

for artifact in "${artifacts[@]}"; do
  if [[ -d "${artifact}" || "${artifact}" == *.asc ]]; then
    continue
  fi
  echo "→ Signing ${artifact}"
  gpg --batch --yes --armor --local-user "${SIGNING_KEY_ID}" --output "${artifact}.asc" --detach-sign "${artifact}"
done

echo "📤 Exporting public key to ${EXPORT_DIR}/public-key.asc"
gpg --armor --export "${SIGNING_KEY_ID}" > "${EXPORT_DIR}/public-key.asc"

echo "✅ Signing complete. Attach .asc files and public-key.asc with release artifacts."
