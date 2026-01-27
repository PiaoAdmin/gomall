#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

pids=()

clean_logs() {
  rm -f "${ROOT_DIR}"/tmp_*.log
  for svc in user product cart order checkout payment api; do
    rm -f "${ROOT_DIR}/app/${svc}/log/kitex.log"
  done
}

if [[ "${1:-}" == "clean" ]]; then
  clean_logs
  echo "Logs cleaned."
  exit 0
fi

clean_logs

start_service() {
  local name="$1"
  local dir="$2"
  echo "Starting ${name}..."
  (cd "$dir" && mkdir -p log && GO_ENV=test go run . >"${ROOT_DIR}/tmp_${name}.log" 2>&1) &
  pids+=("$!")
}

cleanup() {
  echo "Stopping services..."
  for pid in "${pids[@]:-}"; do
    if kill -0 "$pid" 2>/dev/null; then
      kill "$pid" 2>/dev/null || true
    fi
  done
}
trap cleanup EXIT

start_service "user" "${ROOT_DIR}/app/user"
start_service "product" "${ROOT_DIR}/app/product"
start_service "cart" "${ROOT_DIR}/app/cart"
start_service "order" "${ROOT_DIR}/app/order"
start_service "checkout" "${ROOT_DIR}/app/checkout"
start_service "payment" "${ROOT_DIR}/app/payment"

echo "Waiting for services to boot..."
sleep 6

cd "${ROOT_DIR}/app/api"
GO_ENV=test go test ./test -count=1
