#!/usr/bin/env bash

set -euo pipefail

EXPECTED_VERSION="v1.36.2"
EXPECTED_NODE_COUNT="3"
EXPECTED_CONTROL_PLANE_COUNT="1"
EXPECTED_NANO_CPUS="2000000000"
EXPECTED_MEMORY_BYTES="2147483648"
CLUSTER_NAME="${CLUSTER_NAME:-interview-k8s}"
CONTEXT="kind-${CLUSTER_NAME}"

if ! command -v kubectl >/dev/null 2>&1; then
  echo "Missing required command: kubectl" >&2
  exit 1
fi

if ! command -v kind >/dev/null 2>&1 || ! command -v docker >/dev/null 2>&1; then
  echo "Missing required command: kind or docker" >&2
  exit 1
fi

if ! kubectl config get-contexts "${CONTEXT}" >/dev/null 2>&1; then
  echo "Kubernetes context '${CONTEXT}' does not exist. Run 'make cluster-up' first." >&2
  exit 1
fi

echo "Waiting for every node to become Ready..."
kubectl --context "${CONTEXT}" wait --for=condition=Ready nodes --all --timeout=120s

node_count="$(kubectl --context "${CONTEXT}" get nodes --no-headers | wc -l | tr -d ' ')"
control_plane_count="$(
  kubectl --context "${CONTEXT}" get nodes \
    --selector=node-role.kubernetes.io/control-plane \
    --no-headers \
    | wc -l \
    | tr -d ' '
)"

if [[ "${node_count}" != "${EXPECTED_NODE_COUNT}" ]] || [[ "${control_plane_count}" != "${EXPECTED_CONTROL_PLANE_COUNT}" ]]; then
  echo "Expected one control-plane and two worker nodes, but found ${control_plane_count} control-plane node(s) out of ${node_count} total." >&2
  exit 1
fi

while IFS= read -r node; do
  nano_cpus="$(docker inspect --format '{{.HostConfig.NanoCpus}}' "${node}")"
  memory_bytes="$(docker inspect --format '{{.HostConfig.Memory}}' "${node}")"
  memory_swap_bytes="$(docker inspect --format '{{.HostConfig.MemorySwap}}' "${node}")"

  if [[ "${nano_cpus}" != "${EXPECTED_NANO_CPUS}" ]] \
    || [[ "${memory_bytes}" != "${EXPECTED_MEMORY_BYTES}" ]] \
    || [[ "${memory_swap_bytes}" != "${EXPECTED_MEMORY_BYTES}" ]]; then
    echo "Node '${node}' does not have the expected 2 CPU / 2 GiB hard limit." >&2
    exit 1
  fi
done < <(kind get nodes --name "${CLUSTER_NAME}")

server_version="$(
  kubectl --context "${CONTEXT}" get --raw=/version \
    | sed -n 's/.*"gitVersion"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p'
)"

if [[ "${server_version}" != "${EXPECTED_VERSION}" ]]; then
  echo "Expected Kubernetes ${EXPECTED_VERSION}, but the API server reports ${server_version:-unknown}." >&2
  exit 1
fi

echo "Kubernetes API server version: ${server_version}"
echo "Topology: one control-plane and two workers"
echo "Docker limit per node: 2 CPUs and 2 GiB memory (swap disabled)"
kubectl --context "${CONTEXT}" get nodes -o wide
kubectl --context "${CONTEXT}" get pods --all-namespaces
echo "Cluster '${CLUSTER_NAME}' is ready."
