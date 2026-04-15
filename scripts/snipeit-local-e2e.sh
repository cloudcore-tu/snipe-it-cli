#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT_DIR}"

SNIPEIT_URL="${SNIPEIT_URL:-http://localhost:18080}"
SITE_NAME="${SNIPEIT_SITE_NAME:-Local Snipe-IT}"
ADMIN_FIRST_NAME="${SNIPEIT_ADMIN_FIRST_NAME:-Local}"
ADMIN_LAST_NAME="${SNIPEIT_ADMIN_LAST_NAME:-Admin}"
ADMIN_EMAIL="${SNIPEIT_ADMIN_EMAIL:-admin@example.com}"
ADMIN_USERNAME="${SNIPEIT_ADMIN_USERNAME:-admin}"
ADMIN_PASSWORD="${SNIPEIT_ADMIN_PASSWORD:-password}"
TOKEN_NAME="${SNIPEIT_TOKEN_NAME:-local-smoke}"
GOCACHE_DIR="${GOCACHE:-/tmp/snipeit-cli-gocache}"

cleanup() {
  docker compose down -v >/dev/null 2>&1 || true
}

wait_for_http() {
  local attempts=0

  until curl -fsS "${SNIPEIT_URL}/setup" >/dev/null 2>&1; do
    attempts=$((attempts + 1))
    if (( attempts >= 90 )); then
      echo "Snipe-IT did not become ready at ${SNIPEIT_URL}" >&2
      docker compose logs --tail=200 snipeit >&2 || true
      return 1
    fi
    sleep 2
  done
}

extract_csrf_token() {
  sed -n 's/.*name="_token" value="\([^"]*\)".*/\1/p' | head -n1
}

run_setup_wizard() {
  local cookie_file
  local setup_page
  local setup_token
  local user_page
  local user_token

  cookie_file="$(mktemp)"

  setup_page="$(curl -fsS -c "${cookie_file}" "${SNIPEIT_URL}/setup")"
  setup_token="$(printf '%s' "${setup_page}" | extract_csrf_token)"

  curl -fsS -b "${cookie_file}" -c "${cookie_file}" \
    -X POST \
    -d "_token=${setup_token}" \
    "${SNIPEIT_URL}/setup/migrate" >/dev/null

  user_page="$(curl -fsS -b "${cookie_file}" -c "${cookie_file}" "${SNIPEIT_URL}/setup/user")"
  user_token="$(printf '%s' "${user_page}" | extract_csrf_token)"

  curl -fsS -b "${cookie_file}" -c "${cookie_file}" \
    -X POST \
    -d "_token=${user_token}" \
    --data-urlencode "site_name=${SITE_NAME}" \
    --data-urlencode "first_name=${ADMIN_FIRST_NAME}" \
    --data-urlencode "last_name=${ADMIN_LAST_NAME}" \
    --data-urlencode "email=${ADMIN_EMAIL}" \
    --data-urlencode "username=${ADMIN_USERNAME}" \
    --data-urlencode "password=${ADMIN_PASSWORD}" \
    --data-urlencode "password_confirmation=${ADMIN_PASSWORD}" \
    -d "default_currency=USD" \
    -d "locale=en-US" \
    "${SNIPEIT_URL}/setup/user" >/dev/null

  rm -f "${cookie_file}"
}

lookup_admin_user_id() {
  docker compose exec -T db mariadb -N -B -usnipeit -psnipeit -Dsnipeit \
    -e "select id from users where username='${ADMIN_USERNAME}' limit 1;"
}

ensure_personal_access_client() {
  local client_id

  docker compose exec -T snipeit php artisan passport:keys --force >/dev/null

  client_id="$(
    docker compose exec -T db mariadb -N -B -usnipeit -psnipeit -Dsnipeit \
      -e "select id from oauth_clients where personal_access_client = 1 order by id asc limit 1;"
  )"

  if [[ -n "${client_id}" ]]; then
    return 0
  fi

  docker compose exec -T snipeit php artisan passport:client \
    --personal \
    --name="${TOKEN_NAME}" \
    --no-interaction >/dev/null
}

main() {
  local user_id
  local token

  trap cleanup EXIT

  echo "Resetting local Snipe-IT stack"
  docker compose down -v >/dev/null 2>&1 || true

  echo "Starting local Snipe-IT"
  docker compose up -d >/dev/null

  echo "Waiting for setup page"
  wait_for_http

  echo "Running setup wizard"
  run_setup_wizard

  echo "Looking up admin user"
  user_id="$(lookup_admin_user_id)"

  echo "Ensuring Passport keys and personal access client"
  ensure_personal_access_client

  echo "Generating local API token"
  token="$(
    docker compose exec -T snipeit php artisan snipeit:make-api-key \
      --user_id="${user_id}" \
      --name="${TOKEN_NAME}" \
      --key-only
  )"

  echo "Running smoke test"
  env GOCACHE="${GOCACHE_DIR}" SNIPEIT_URL="${SNIPEIT_URL}" SNIPEIT_TOKEN="${token}" \
    bash scripts/snipeit-local-smoke.sh

  echo "Local end-to-end smoke test passed"
}

main "$@"
