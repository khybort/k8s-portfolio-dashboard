#!/usr/bin/env sh
set -eu
# Enable pipefail when available (zsh/bash), ignore if unsupported (dash/sh).
(set -o pipefail) >/dev/null 2>&1 && set -o pipefail || true

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
IMAGE="${IMAGE:-portfolio-dashboard:local}"
NAMESPACE="${NAMESPACE:-portfolio}"

need() { command -v "$1" >/dev/null 2>&1 || { echo "Missing dependency: $1" >&2; exit 1; }; }

need kubectl

CTX="$(kubectl config current-context 2>/dev/null || true)"
echo "Using kubectl context: ${CTX:-<unknown>}"
echo "Image: ${IMAGE}"

case "${CTX}" in
  *kind*)
  need docker
  need kind
  echo "Detected kind. Building image with docker and loading into kind node..."
  docker build -t "${IMAGE}" "${ROOT_DIR}"
  kind load docker-image "${IMAGE}"
  ;;
  *minikube*)
  need minikube
  echo "Detected minikube. Building image inside minikube..."
  # avoids eval'ing docker-env and works across shells
  minikube image build -t "${IMAGE}" "${ROOT_DIR}"
  ;;
  *)
  # Docker Desktop Kubernetes / Colima (shared daemon) or already-available image
  if command -v docker >/dev/null 2>&1; then
    echo "Building image with docker (context is not kind/minikube)..."
    docker build -t "${IMAGE}" "${ROOT_DIR}"
  else
    echo "docker not found. Assuming image '${IMAGE}' is already available to the cluster runtime."
  fi
  ;;
esac

echo "Applying manifests..."
kubectl apply -k "${ROOT_DIR}/k8s"
kubectl -n "${NAMESPACE}" rollout status deploy/portfolio-dashboard

echo
echo "Deployed. To access locally:"
echo "  ${ROOT_DIR}/port-forward.sh"


