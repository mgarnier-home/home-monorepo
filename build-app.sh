#!/bin/bash
set -euo pipefail

# Default values
APP_NAME=""
VERSION="test"
TAG="latest"
USE_CACHE="true"
PROGRESS="false"

# Function to display usage
usage() {
  echo "Usage: $0 --name <app_name> [--version <VERSION>] [--tag <tag>] [--no-cache] [--progress]"
  exit 1
}

# Parse named parameters
while [[ "$#" -gt 0 ]]; do
  case "$1" in
    --name)
      APP_NAME="$2"
      shift 2
      ;;
    --version)
      VERSION="$2"
      shift 2
      ;;
    --tag)
      TAG="$2"
      shift 2
      ;;
    --no-cache)
      USE_CACHE="false"
      shift
      ;;
    --progress)
      PROGRESS="true"
      shift
      ;;
    *)
      echo "Unknown parameter: $1"
      usage
      ;;
  esac
done

# Check if APP_NAME is provided
if [ -z "${APP_NAME}" ]; then
  usage
fi

echo "Building app : $APP_NAME version : $VERSION tag : $TAG"

docker rmi build-node-app || true
docker rmi build-go-app || true
echo "Deleted build images"

BUILD_IMAGE_ARGS=("--build-arg" "APP=$APP_NAME" "--build-arg" "VERSION=$VERSION" ".")
RUNTIME_IMAGE_ARGS=("-t" "mgarnier11/$APP_NAME:$TAG" "-f" "apps/$APP_NAME/docker/Dockerfile" ".")


if [[ "$USE_CACHE" == "false" ]]; then
  BUILD_IMAGE_ARGS=("--no-cache" "${BUILD_IMAGE_ARGS[@]}")
  RUNTIME_IMAGE_ARGS=("--no-cache" "${RUNTIME_IMAGE_ARGS[@]}")
fi

if [[ "$PROGRESS" == "true" ]]; then
  BUILD_IMAGE_ARGS=("--progress" "plain" "${BUILD_IMAGE_ARGS[@]}")
  RUNTIME_IMAGE_ARGS=("--progress" "plain" "${RUNTIME_IMAGE_ARGS[@]}")
fi

if [[ -f "apps/$APP_NAME/package.json" ]]; then
  echo "Node app detected"
  docker build -t build-node-app -f docker/node.build.Dockerfile "${BUILD_IMAGE_ARGS[@]}"
fi
if [[ -f "apps/$APP_NAME/taskfile.yml" ]]; then
  echo "Golang app detected"
  docker build -t build-go-app -f docker/golang.build.Dockerfile "${BUILD_IMAGE_ARGS[@]}"
fi

echo "Built new build images for $APP_NAME"

docker build "${RUNTIME_IMAGE_ARGS[@]}"
echo "Built new runtime image for $APP_NAME"

docker rmi build-node-app || true
docker rmi build-go-app || true
echo "Cleaned up build images"
