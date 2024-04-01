#!/bin/bash

PROJECT_NAME=${1}
BUILDX=${2: false}
PUSH=${3: false}

if [ -z "${PROJECT_NAME}" ]; then
  echo "Usage: $0 <project_name>"
  exit 1
fi

if [ "${BUILDX}" = true ]; then
  echo "Building for multiple platforms"
  docker buildx create --name mybuilder
  docker buildx use mybuilder
  docker buildx --no-cache build --platform linux/amd64,linux/arm64 -t mgarnier11/$PROJECT_NAME:latest -f apps/$PROJECT_NAME/docker/Dockerfile .
else
  echo "Building for single platform"
  docker build --no-cache -t mgarnier11/$PROJECT_NAME:latest -f apps/$PROJECT_NAME/docker/Dockerfile .
fi

if [ "${PUSH}" = true ]; then
  echo "Pushing to Docker Hub"
  docker push mgarnier11/$PROJECT_NAME:latest
fi

