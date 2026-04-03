#!/usr/bin/env bash
set -euo pipefail

IMAGE_NAME="${IMAGE_NAME:-torq/ip2country-service}"
IMAGE_TAG="${IMAGE_TAG:-latest}"
DOCKERFILE_PATH="${DOCKERFILE_PATH:-Dockerfile}"
BUILD_CONTEXT="${BUILD_CONTEXT:-.}"
NO_CACHE="${NO_CACHE:-false}"

BUILD_ARGS=()
if [ "${NO_CACHE}" = "true" ]; then
  BUILD_ARGS+=(--no-cache)
fi

docker build "${BUILD_ARGS[@]}" -f "${DOCKERFILE_PATH}" -t "${IMAGE_NAME}:${IMAGE_TAG}" "${BUILD_CONTEXT}"
echo "Built ${IMAGE_NAME}:${IMAGE_TAG}"
