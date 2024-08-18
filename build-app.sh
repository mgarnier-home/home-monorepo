APP_NAME=${1:-""}

APP_VERSION=${2:-"test"}

USE_CACHE=${3:-"true"}
PROGRESS=${4:-"false"}

if [ -z "${APP_NAME}" ]; then
  echo "Usage: $0 <app_name>"
  exit 1
fi

echo "Building app : $APP_NAME version : $APP_VERSION"


docker rmi build
echo "Deleted build image"

BUILD_IMAGE_ARGS=("-t" "build" "--build-arg" "APP=$APP_NAME" "--build-arg" "APP_VERSION=$APP_VERSION" ".")
RUNTIME_IMAGE_ARGS=("-t" "mgarnier11/$APP_NAME:latest" "-f" "apps/$APP_NAME/docker/Dockerfile" ".")
# BUILD_IMAGE_ARGS="-t build -f docker/build.Dockerfile --build-arg APP=$PROJECT_NAME --build-arg APP_VERSION=test --progress plain ."
# DOCKER_BUILD_RUN_ARGS="--no-cache -t mgarnier11/$PROJECT_NAME:latest -f apps/$PROJECT_NAME/docker/Dockerfile --progress plain ."


if [[ -f "apps/$APP_NAME/package.json" ]]; then
  echo "Node app detected"
  BUILD_IMAGE_ARGS=("-f" "docker/node.build.Dockerfile" "${BUILD_IMAGE_ARGS[@]}")
elif [[ -f "apps/$APP_NAME/go.mod" ]]; then
  echo "Golang app detected"
  BUILD_IMAGE_ARGS=("-f" "docker/golang.build.Dockerfile" "${BUILD_IMAGE_ARGS[@]}")
fi

if [[ "$USE_CACHE" == "false" ]]; then
  BUILD_IMAGE_ARGS=("--no-cache" "${BUILD_IMAGE_ARGS[@]}")
  RUNTIME_IMAGE_ARGS=("--no-cache" "${RUNTIME_IMAGE_ARGS[@]}")
fi

if [[ "$PROGRESS" == "true" ]]; then
  BUILD_IMAGE_ARGS=("--progress" "plain" "${BUILD_IMAGE_ARGS[@]}")
  RUNTIME_IMAGE_ARGS=("--progress" "plain" "${RUNTIME_IMAGE_ARGS[@]}")
fi

docker build "${BUILD_IMAGE_ARGS[@]}"
echo "Built new build image for $APP_NAME"

docker build "${RUNTIME_IMAGE_ARGS[@]}"
echo "Built new runtime image for $APP_NAME"

docker rmi build
echo "Cleaned up build image"
