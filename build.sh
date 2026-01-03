#!/usr/bin/env sh
set -euo pipefail

# Build script for Portfolio Dashboard
# Usage: ./build.sh [tag]

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
IMAGE_TAG="${1:-local}"
IMAGE_NAME="portfolio-dashboard"

echo "Building Docker image: ${IMAGE_NAME}:${IMAGE_TAG}"

cd "${ROOT_DIR}"

# Build image
docker build \
  --tag "${IMAGE_NAME}:${IMAGE_TAG}" \
  --tag "${IMAGE_NAME}:latest" \
  .

echo "✅ Build complete: ${IMAGE_NAME}:${IMAGE_TAG}"

# Show image info
docker images "${IMAGE_NAME}:${IMAGE_TAG}" --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"

