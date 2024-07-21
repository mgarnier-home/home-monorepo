#!/bin/bash

PROJECT_NAME=${1}
BUILDX=${2: false}
PUSH=${3: false}

if [ -z "${PROJECT_NAME}" ]; then
  echo "Usage: $0 <project_name>"
  exit 1
fi

DOCKER_BUILD_APP_ARGS="--no-cache -t build -f docker/build.Dockerfile --build-arg APP=$PROJECT_NAME --build-arg APP_VERSION=test ."
DOCKER_BUILD_RUN_ARGS="--no-cache -t mgarnier11/$PROJECT_NAME:latest -f apps/$PROJECT_NAME/docker/Dockerfile --progress plain ."

if [ "${BUILDX}" = true ]; then
  echo "Building for multiple platforms"
  docker buildx create --name mybuilder
  docker buildx use mybuilder
  docker buildx build --platform linux/amd64,linux/arm64 $DOCKER_BUILD_APP_ARGS
  docker buildx build --platform linux/amd64,linux/arm64 $DOCKER_BUILD_RUN_ARGS
else
  echo "Building for single platform"
  docker build $DOCKER_BUILD_APP_ARGS
  docker build $DOCKER_BUILD_RUN_ARGS
fi

if [ "${PUSH}" = true ]; then
  echo "Pushing to Docker Hub"
  docker push mgarnier11/$PROJECT_NAME:latest
fi

