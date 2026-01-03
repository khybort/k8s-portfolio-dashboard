#!/usr/bin/env sh
set -eu
(set -o pipefail) >/dev/null 2>&1 && set -o pipefail || true

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
command -v kubectl >/dev/null 2>&1 || { echo "Missing dependency: kubectl" >&2; exit 1; }

kubectl delete -k "${ROOT_DIR}/k8s"


