#!/usr/bin/env bash

set -euo pipefail

CLUSTER_NAME="${CLUSTER_NAME:-interview-k8s}"

if ! command -v kind >/dev/null 2>&1; then
  echo "Missing required command: kind" >&2
  exit 1
fi

if kind get clusters 2>/dev/null | grep -Fxq "${CLUSTER_NAME}"; then
  kind delete cluster --name "${CLUSTER_NAME}"
else
  echo "Cluster '${CLUSTER_NAME}' does not exist; nothing to delete."
fi

