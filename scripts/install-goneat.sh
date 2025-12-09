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









GONEAT_VERSION="v0.3.14"









# Map uname outputs to archive naming




OS_RAW="$(uname -s)"




case "${OS_RAW}" in




  Darwin*) OS="darwin" ;;




  Linux*) OS="linux" ;;




  MINGW* | MSYS* | CYGWIN*) OS="windows" ;;




  *) OS="$(echo "${OS_RAW}" | tr '[:upper:]' '[:lower:]')" ;;




esac









ARCH_RAW="$(uname -m)"




case "${ARCH_RAW}" in




  x86_64) ARCH="amd64" ;;




  aarch64 | arm64) ARCH="arm64" ;;




  *) echo "Unsupported architecture: ${ARCH_RAW}" >&2; exit 1 ;;




esac









# Windows uses .zip, others use .tar.gz




if [[ "${OS}" == "windows" ]]; then




  ARCHIVE="goneat_${GONEAT_VERSION}_${OS}_${ARCH}.zip"




else




  ARCHIVE="goneat_${GONEAT_VERSION}_${OS}_${ARCH}.tar.gz"




fi




URL="https://github.com/fulmenhq/goneat/releases/download/${GONEAT_VERSION}/${ARCHIVE}"




case "${OS}-${ARCH}" in




  darwin-amd64)
    EXPECTED_SHA="25ed351c4733670548d6de5eb5d9d284565d45111dee0180787f9bbd7b5f0a9a"
    ;;




  darwin-arm64)
    EXPECTED_SHA="c5cbe99c559a3fb1c38b0b8a65780493131b64757db827be6786277d18e85620"
    ;;




  linux-amd64)
    EXPECTED_SHA="b9c0020824a344a92179a5130b8ce3b7aed516c9bfa054a1fb809d41176e7349"
    ;;




  linux-arm64)
    EXPECTED_SHA="07375fb31cb94b4821b7b3d08cb03529d1d1f12e4b40b371edd032ab56077442"
    ;;




  windows-amd64)
    EXPECTED_SHA="e7f442ed05937602523be8950b497d0294f97fe829cd2acf8c61f73c64408cc9"
    ;;




  windows-arm64)




    EXPECTED_SHA="PENDING"




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




    # Use sha256sum on Linux, shasum on macOS
    if command -v sha256sum >/dev/null 2>&1; then
      ACTUAL_SHA="$(sha256sum "${tmpfile}" | cut -d' ' -f1)"
    else
      ACTUAL_SHA="$(shasum -a 256 "${tmpfile}" | cut -d' ' -f1)"
    fi




    if [[ "${ACTUAL_SHA}" != "${EXPECTED_SHA}" ]]; then




      echo "❌ Checksum mismatch!" >&2




      echo "   Expected: ${EXPECTED_SHA}" >&2




      echo "   Actual:   ${ACTUAL_SHA}" >&2




      exit 1




    fi




    echo "✅ Checksum verified"




  fi









  # Extract archive (zip for Windows, tar.gz for others)




  if [[ "${OS}" == "windows" ]]; then




    unzip -q -o "${tmpfile}" -d "${BIN_DIR}"




    # Windows binary may be goneat.exe, rename if needed




    if [[ -f "${BIN_DIR}/goneat.exe" ]]; then




      mv "${BIN_DIR}/goneat.exe" "${BIN_DIR}/goneat"




    fi




  else




    tar -xzf "${tmpfile}" -C "${BIN_DIR}"




  fi




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
"${BIN_DIR}/goneat" doctor tools --scope foundation --install --install-package-managers --yes --no-cooling









echo "✅ Bootstrap complete"
