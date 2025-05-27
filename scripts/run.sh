#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
readonly PROJECT_ROOT_DIR
readonly BINARY_PATH="${PROJECT_ROOT_DIR}/bin/go_api_server"
readonly LOG_DIR="${PROJECT_ROOT_DIR}/logs"
readonly LOG_FILE="${LOG_DIR}/app.log"

cd "${PROJECT_ROOT_DIR}"

echo "[INFO] Starting application in background..."

echo "[INFO] Checking application binary: ${BINARY_PATH}"
if [[ ! -x "${BINARY_PATH}" ]]; then
    echo "[ERROR] Binary not found or not executable: ${BINARY_PATH}" >&2
    exit 1
fi

echo "[INFO] Ensuring log directory exists: ${LOG_DIR}"
if ! mkdir -p "${LOG_DIR}"; then
    echo "[ERROR] Failed to create log directory: ${LOG_DIR}" >&2
    exit 1
fi

echo "[INFO] Starting application with nohup, logging to: ${LOG_FILE}"
if nohup "${BINARY_PATH}" "$@" >"${LOG_FILE}" 2>&1 & then
    pid=$!
    echo "[INFO] Application started in background. PID (of nohup): ${pid}"
else
    echo "[ERROR] Failed to start application with nohup." >&2
    exit 1
fi

exit 0
