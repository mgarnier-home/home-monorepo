#!/bin/bash

PROJECT_NAME=${1}

if [ -z "${PROJECT_NAME}" ]; then
  echo "Usage: $0 <project_name>"
  exit 1
fi

docker buildx build --platform linux/amd64,linux/arm64 -t mgarnier11/$PROJECT_NAME:latest -f apps/$PROJECT_NAME/docker/Dockerfile . --push
