#!/bin/bash
# set -e

PROJECT_NAME=${1}


if [ -z "${PROJECT_NAME}" ]; then
  echo "Usage: $0 <project_name>"
  exit 1
fi

echo "Building $PROJECT_NAME"

DOCKER_BUILD_APP_ARGS="--no-cache -t build -f docker/build.Dockerfile --build-arg APP=$PROJECT_NAME --build-arg APP_VERSION=test --progress plain ."
DOCKER_BUILD_RUN_ARGS="--no-cache -t mgarnier11/$PROJECT_NAME:latest -f apps/$PROJECT_NAME/docker/Dockerfile --progress plain ."

echo "Building for single platform"
docker rmi build
docker build $DOCKER_BUILD_APP_ARGS

docker build $DOCKER_BUILD_RUN_ARGS
