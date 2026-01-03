#!/usr/bin/env sh
set -eu
(set -o pipefail) >/dev/null 2>&1 && set -o pipefail || true

NAMESPACE="${NAMESPACE:-portfolio}"

command -v kubectl >/dev/null 2>&1 || { echo "Missing dependency: kubectl" >&2; exit 1; }

echo "Port-forwarding svc/portfolio-dashboard (namespace: ${NAMESPACE}) to http://localhost:8080 ..."
echo "Press Ctrl+C to stop."
kubectl -n "${NAMESPACE}" port-forward svc/portfolio-dashboard 8080:80


