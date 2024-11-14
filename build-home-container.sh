#!/bin/bash

# Default values
VERSION="test"
TAG="latest"
USE_CACHE="true"
PROGRESS="false"

# Function to display usage
usage() {
  echo "Usage: $0 [--version <VERSION>] [--tag <tag>] [--no-cache] [--progress]"
  exit 1
}


# Parse named parameters
while [[ "$#" -gt 0 ]]; do
  case "$1" in
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

echo "Building home-container : version : $VERSION tag : $TAG"

BUILD_ARGS=("-t" "mgarnier11/home-container:$TAG" "-t" "mgarnier11/home-container:$VERSION" "--build-arg" "VERSION=$VERSION" "-f" ".devcontainer/Dockerfile" "./.devcontainer")

if [[ "$USE_CACHE" == "false" ]]; then
  BUILD_ARGS=("--no-cache" "${BUILD_ARGS[@]}")
fi

if [[ "$PROGRESS" == "true" ]]; then
  BUILD_ARGS=("--progress" "plain" "${BUILD_ARGS[@]}")
fi

docker build "${BUILD_ARGS[@]}"
echo "Built home-container image"

