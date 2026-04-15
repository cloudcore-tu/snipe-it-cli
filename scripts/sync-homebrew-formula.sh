#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 3 ]]; then
  echo "usage: $0 TAG CHECKSUMS_PATH TAP_DIR" >&2
  exit 1
fi

tag="$1"
checksums_path="$2"
tap_dir="$3"
version="${tag#v}"
formula_path="${tap_dir}/Formula/snipe-it-cli.rb"
push_changes="${SYNC_HOMEBREW_PUSH:-1}"

if [[ -z "${version}" ]]; then
  echo "TAG must not be empty" >&2
  exit 1
fi

if [[ ! -f "${checksums_path}" ]]; then
  echo "checksums file not found: ${checksums_path}" >&2
  exit 1
fi

if [[ ! -d "${tap_dir}" ]]; then
  echo "tap directory not found: ${tap_dir}" >&2
  exit 1
fi

if [[ ! -f "${tap_dir}/scripts/update-snipe-it-cli-formula.py" ]]; then
  echo "formula update script not found in tap repo: ${tap_dir}/scripts/update-snipe-it-cli-formula.py" >&2
  exit 1
fi

python3 "${tap_dir}/scripts/update-snipe-it-cli-formula.py" "${version}" "${checksums_path}"

cd "${tap_dir}"
git config user.email "github-actions[bot]@users.noreply.github.com"
git config user.name "github-actions[bot]"
git add "${formula_path}"

if git diff --cached --quiet; then
  echo "No formula changes"
  exit 0
fi

git commit -m "chore: snipe-it-cli formula を v${version} に更新"

if [[ "${push_changes}" != "1" ]]; then
  echo "Skipping git push because SYNC_HOMEBREW_PUSH=${push_changes}"
  exit 0
fi

git push
