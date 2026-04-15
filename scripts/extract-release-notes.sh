#!/usr/bin/env bash

set -euo pipefail

if [[ $# -lt 2 || $# -gt 3 ]]; then
  echo "usage: $0 VERSION OUTPUT_PATH [CHANGELOG_PATH]" >&2
  exit 1
fi

version="$1"
output_path="$2"
changelog_path="${3:-CHANGELOG.md}"

if [[ -z "${version}" ]]; then
  echo "VERSION must not be empty" >&2
  exit 1
fi

if [[ ! -f "${changelog_path}" ]]; then
  echo "CHANGELOG file not found: ${changelog_path}" >&2
  exit 1
fi

tmp_output="$(mktemp)"
trap 'rm -f "${tmp_output}"' EXIT

awk -v version="${version}" '
  $0 ~ "^## \\[" version "\\]" { flag=1; next }
  /^## \[/ && flag && $0 !~ "^## \\[" version "\\]" { exit }
  flag { print }
' "${changelog_path}" \
  | sed '/^[[:space:]]*$/d' \
  > "${tmp_output}"

if [[ ! -s "${tmp_output}" ]]; then
  echo "release notes section for v${version} not found or empty in ${changelog_path}" >&2
  exit 1
fi

mkdir -p "$(dirname "${output_path}")"
cp "${tmp_output}" "${output_path}"

echo "Release notes for v${version}:"
cat "${output_path}"
