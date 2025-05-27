#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
readonly PROJECT_ROOT_DIR

readonly MAIN_GO_FILE="cmd/api/main.go"
readonly OUTPUT_BINARY="${PROJECT_ROOT_DIR}/bin/go_api_server"
readonly OUTPUT_DIR="$(dirname "${OUTPUT_BINARY}")"

cd "${PROJECT_ROOT_DIR}"

echo "[INFO] Generating Swagger documentation..."
if ! swag init -g "${MAIN_GO_FILE}"; then
    echo "[ERROR] Swagger documentation generation failed. Check output." >&2
    exit 1
fi

echo "[INFO] Starting Go build process..."

echo "[INFO] Ensuring output directory: ${OUTPUT_DIR}"
if ! mkdir -p "${OUTPUT_DIR}"; then
    echo "[ERROR] Failed to create output directory: ${OUTPUT_DIR}" >&2
    exit 1
fi

echo "[INFO] Building Go binary..."
if ! go build \
    -a \
    -trimpath \
    -buildvcs=false \
    -ldflags="-s -w" \
    -installsuffix cgo \
    -o "${OUTPUT_BINARY}" \
    "${MAIN_GO_FILE}"; then
    echo "[ERROR] Go build failed for '${MAIN_GO_FILE}'. Check output." >&2
    exit 1
fi

echo "[INFO] Setting execute permissions..."
if ! chmod +x "${OUTPUT_BINARY}"; then
    echo "[ERROR] Failed to set execute permissions on ${OUTPUT_BINARY}." >&2
    exit 1
fi

echo "[INFO] Build successful!"
echo "[INFO] Binary available at: ${OUTPUT_BINARY}"
exit 0
