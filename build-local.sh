#!/bin/bash

PROJECT_NAME=${1}

# BUILDX=${2: false}
# PUSH=${3: false}

if [ -z "${PROJECT_NAME}" ]; then
  echo "Usage: $0 <project_name>"
  exit 1
fi

echo "Building $PROJECT_NAME"
echo "New build: $NEW_BUILD"

DOCKER_BUILD_APP_ARGS="--no-cache -t build -f docker/build.Dockerfile --build-arg APP=$PROJECT_NAME --build-arg APP_VERSION=test --progress plain ."
DOCKER_BUILD_RUN_ARGS="--no-cache -t mgarnier11/$PROJECT_NAME:latest -f apps/$PROJECT_NAME/docker/Dockerfile --progress plain ."

# if [ "${BUILDX}" = true ]; then
#   echo "Building for multiple platforms"
#   # docker buildx create --name mybuilder
#   docker buildx use mybuilder
#   docker buildx build --platform linux/amd64,linux/arm64 $DOCKER_BUILD_APP_ARGS --load
#   docker buildx build --platform linux/amd64,linux/arm64 $DOCKER_BUILD_RUN_ARGS
# else
echo "Building for single platform"
docker rmi build
docker build $DOCKER_BUILD_APP_ARGS

docker build $DOCKER_BUILD_RUN_ARGS
# fi

# if [ "${PUSH}" = true ]; then
#   echo "Pushing to Docker Hub"
#   docker push mgarnier11/$PROJECT_NAME:latest
# fi



# docker buildx build --platform linux/arm64 --no-cache -t build:arm64 -f docker/build.Dockerfile --build-arg APP=autosaver --build-arg APP_VERSION=test --load .
# docker buildx build --platform linux/arm64 --no-cache -t mgarnier11/autosaver:arm64 -f apps/autosaver/docker/Dockerfile . 
# docker build --no-cache -t mgarnier11/autosaver:arm64 -f apps/autosaver/docker/Dockerfile . 