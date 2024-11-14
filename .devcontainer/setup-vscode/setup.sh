#!/usr/bin/env bash

set -euo pipefail

get_value_with_default() {
  local json_input=$1
  local key=$2
  local default_value=$3
  local result

  result=$(jq -r --argjson data "$json_input" --arg key "$key" --arg default "$default_value" -n '($data[$key] | select(length > 0) // $default)')

  echo "$result"
}

echo "First argument is $1"

INSTALL_PREFIX=$(get_value_with_default "$1" "installPrefix" "/tmp/vscode-web")
EXTENSIONS_DIR=$(get_value_with_default "$1" "extensionsDir" "/home/${USER}/.vscode-web/extensions")
EXTENSIONS=$(get_value_with_default "$1" "extensions" "[]")
LOG_PATH=$(get_value_with_default "$1" "logPath" "/tmp/vscode-web.log")
PORT=$(get_value_with_default "$1" "port" "13338")

echo "INSTALL_PREFIX is $INSTALL_PREFIX"
echo "EXTENSIONS_DIR is $EXTENSIONS_DIR"
echo "EXTENSIONS is $EXTENSIONS"

VSCODE_WEB="$INSTALL_PREFIX/bin/code-server"

echo "Installing Microsoft Visual Studio Code Server! to ${INSTALL_PREFIX}"

# Create install prefix
mkdir -p ${INSTALL_PREFIX}

# Download and extract vscode-server
ARCH=$(uname -m)
case "$ARCH" in
x86_64) ARCH="x64" ;;
aarch64) ARCH="arm64" ;;
*)
  echo "Unsupported architecture"
  exit 1
  ;;
esac

HASH=$(curl -fsSL https://update.code.visualstudio.com/api/commits/stable/server-linux-$ARCH-web | cut -d '"' -f 2)
output=$(curl -fsSL https://vscode.download.prss.microsoft.com/dbazure/download/stable/$HASH/vscode-server-linux-$ARCH-web.tar.gz | tar -xz -C ${INSTALL_PREFIX} --strip-components 1)

if [ $? -ne 0 ]; then
  echo "Failed to install Microsoft Visual Studio Code Server: $output"
  exit 1
fi
echo "VS Code Web has been installed."

# Set extension directory
if [[ -n "${EXTENSIONS_DIR}" ]]; then
  echo "EXTENSIONS_DIR variable is set to $EXTENSIONS_DIR, creating the dir"
  mkdir -p "$EXTENSIONS_DIR"
else
  echo "EXTENSIONS_DIR variable is not set"
fi

EXTENSION_ARG=""
if [ -n "$EXTENSIONS_DIR" ]; then
  EXTENSION_ARG="--extensions-dir=$EXTENSIONS_DIR"
fi

# Install each extension...
echo "$EXTENSIONS" | jq -r '.[]' | while read extension; do
  if [ -z "$extension" ]; then
    continue
  fi
  echo "Installing extension $extension..."
  output=$($VSCODE_WEB "$EXTENSION_ARG" --install-extension "$extension" --force)
  if [ $? -ne 0 ]; then
    echo "Failed to install extension: $extension: $output"
    exit 1
  fi
done

if [ -f "/entrypoint.sh" ]; then
  sudo sed -i "1a ${VSCODE_WEB} serve-web $EXTENSION_ARG --port $PORT --host 0.0.0.0 --accept-server-license-terms --without-connection-token --telemetry-level off > $LOG_PATH 2>&1 &" /entrypoint.sh
  echo "entrypoint.sh has been updated"
fi
