#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
readonly PROJECT_ROOT_DIR
readonly BINARY_PATH="${PROJECT_ROOT_DIR}/bin/go_api_server"
readonly STOP_TIMEOUT_SECONDS=10

echo "[INFO] Attempting to stop application: ${BINARY_PATH}"

PIDS=$(pgrep -f "${BINARY_PATH}" || true)

if [[ -z "$PIDS" ]]; then
    echo "[INFO] Application is not currently running."
    exit 0
fi

echo "[INFO] Found running application instance(s) with PID: ${PIDS}"
echo "[INFO] Sending SIGTERM to initiate graceful shutdown."

if ! kill ${PIDS} 2>/dev/null; then
    echo "[ERROR] Failed to send SIGTERM to PID: ${PIDS}. They may have already stopped or you lack permissions." >&2
fi

for ((i = 0; i < ${STOP_TIMEOUT_SECONDS}; i++)); do
    if ! pgrep -f "${BINARY_PATH}" >/dev/null; then
        echo "[INFO] Application stopped gracefully."
        exit 0
    fi
    sleep 1
    echo -n "."
done
echo

echo "[WARN] Application did not stop gracefully with SIGTERM after ${STOP_TIMEOUT_SECONDS}s."

PIDS_TO_KILL=$(pgrep -f "${BINARY_PATH}" || true)

if [[ -n "$PIDS_TO_KILL" ]]; then
    echo "[INFO] Sending SIGKILL to remaining PID: ${PIDS_TO_KILL}"
    if kill -9 ${PIDS_TO_KILL} 2>/dev/null; then
        echo "[INFO] Successfully sent SIGKILL to PID: ${PIDS_TO_KILL}"
    else
        if ! pgrep -f "${BINARY_PATH}" >/dev/null; then
            echo "[INFO] Application seems to have stopped just before/during SIGKILL."
        else
            echo "[ERROR] Failed to send SIGKILL to PID: ${PIDS_TO_KILL}. Manual intervention may be required (e.g. permissions)." >&2
            exit 1
        fi
    fi
else
    echo "[INFO] Application seems to have stopped during the final check before SIGKILL."
fi

exit 0
