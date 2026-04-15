#!/usr/bin/env bash

set -euo pipefail

: "${SNIPEIT_TOKEN:?set SNIPEIT_TOKEN to a local Snipe-IT API token first}"

SNIPEIT_URL="${SNIPEIT_URL:-http://localhost:18080}"

echo "Running local smoke test against ${SNIPEIT_URL}"

go run . version
go run . --url "${SNIPEIT_URL}" --token "${SNIPEIT_TOKEN}" settings login-attempts -o json >/dev/null
go run . --url "${SNIPEIT_URL}" --token "${SNIPEIT_TOKEN}" account tokens -o json >/dev/null
go run . --url "${SNIPEIT_URL}" --token "${SNIPEIT_TOKEN}" assets list --limit 1 -o json >/dev/null

echo "Smoke test passed."
