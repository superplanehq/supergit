#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

if [ "${1-}" = "" ]; then
  echo "Usage: scripts/build-image.sh <version> <arch>"
  exit 1
fi

if [ "${2-}" = "" ]; then
  echo "Usage: scripts/build-image.sh <version> <arch>"
  exit 1
fi

VERSION="$1"
ARCH="$2"

IMAGE_REPO="${SUPERGIT_IMAGE_REPO:-ghcr.io/superplanehq/supergit}"

echo "Building supergit image (${IMAGE_REPO}:${VERSION}-${ARCH})"

docker buildx build \
  --platform "linux/${ARCH}" \
  --progress=plain \
  --provenance=false \
  --push \
  -t "${IMAGE_REPO}:${VERSION}-${ARCH}" \
  -f Dockerfile \
  .
