#!/usr/bin/env bash
set -euo pipefail

APP_NAME="server-sfp-sla"
CMD_PATH="./cmd/server-sfp-sla"
BUILD_DIR="build"

echo "[1/3] Creating build directory"
mkdir -p "${BUILD_DIR}"

echo "[2/3] Building ${APP_NAME}"
go build -o "${BUILD_DIR}/${APP_NAME}" "${CMD_PATH}"

echo "[3/3] Done"
echo "Binary: ${BUILD_DIR}/${APP_NAME}"
