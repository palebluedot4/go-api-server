#!/usr/bin/env bash
set -euo pipefail

echo "[INFO] Pre-commit checks starting..."

echo "[INFO] Formatting Go files..."
if ! go fmt ./...; then
    echo "[ERROR] 'go fmt' failed. Check output." >&2
    exit 1
fi

echo "[INFO] Tidying Go modules..."
if ! go mod tidy; then
    echo "[ERROR] 'go mod tidy' failed. Check output." >&2
    exit 1
fi

if command -v golangci-lint &>/dev/null; then
    echo "[INFO] Running linter..."
    if ! golangci-lint run --fast ./...; then
        echo "[ERROR] Linter found issues. Check output." >&2
        exit 1
    fi
else
    echo "[WARN] Linter (golangci-lint) not found, skipping."
fi

readonly MAIN_GO_FILE="cmd/api/main.go"
readonly OUTPUT_DIR="bin"
readonly OUTPUT_BINARY="${OUTPUT_DIR}/go_api_server"

echo "[INFO] Building Go binary (check)..."

if ! mkdir -p "${OUTPUT_DIR}"; then
    echo "[ERROR] Failed to create output dir: ${OUTPUT_DIR}" >&2
    exit 1
fi

if ! go build \
    -a \
    -trimpath \
    -buildvcs=false \
    -ldflags="-s -w" \
    -installsuffix cgo \
    -o "${OUTPUT_BINARY}" \
    "${MAIN_GO_FILE}"; then
    echo "[ERROR] Build failed for '${MAIN_GO_FILE}'. Check output." >&2
    exit 1
fi

echo "[INFO] Pre-commit checks passed."
exit 0
