#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT_DIR="${ROOT_DIR}/.artifacts/release"
COMPLETION_DIR="${OUT_DIR}/completions"
MAN_DIR="${OUT_DIR}/man"
GOCACHE_DIR="${GOCACHE:-/tmp/snipeit-cli-release-gocache}"

rm -rf "${OUT_DIR}"
mkdir -p "${COMPLETION_DIR}" "${MAN_DIR}"

cd "${ROOT_DIR}"

env GOCACHE="${GOCACHE_DIR}" go run . completion bash > "${COMPLETION_DIR}/snip.bash"
env GOCACHE="${GOCACHE_DIR}" go run . completion zsh > "${COMPLETION_DIR}/_snip"
env GOCACHE="${GOCACHE_DIR}" go run . completion fish > "${COMPLETION_DIR}/snip.fish"

cp man/en/snip.1 "${MAN_DIR}/snip.en.1"
cp man/ja/snip.1 "${MAN_DIR}/snip.ja.1"
