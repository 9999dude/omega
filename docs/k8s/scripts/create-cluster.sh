#!/usr/bin/env bash

set -euo pipefail

K8S_VERSION="v1.36.2"
CLUSTER_NAME="${CLUSTER_NAME:-interview-k8s}"
NODE_IMAGE="${NODE_IMAGE:-interview/kind-node:${K8S_VERSION}}"
NODE_CPUS="2"
NODE_MEMORY="2g"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LAB_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
CONTEXT="kind-${CLUSTER_NAME}"
CACHE_DIR="${LAB_DIR}/.tools/cache"

require() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    exit 1
  fi
}

require docker
require kind
require kubectl
require curl
require shasum

if ! docker info >/dev/null 2>&1; then
  echo "Docker is installed but is not running. Start Docker, then run 'make cluster-up' again." >&2
  exit 1
fi

if kind get clusters 2>/dev/null | grep -Fxq "${CLUSTER_NAME}"; then
  echo "Cluster '${CLUSTER_NAME}' already exists; verifying it instead of replacing it."
  "${SCRIPT_DIR}/verify-cluster.sh"
  exit 0
fi

if ! docker image inspect "${NODE_IMAGE}" >/dev/null 2>&1; then
  machine="$(uname -m)"
  case "${machine}" in
    x86_64) arch="amd64" ;;
    arm64|aarch64) arch="arm64" ;;
    *)
      echo "Unsupported CPU architecture: ${machine}" >&2
      exit 1
      ;;
  esac

  mkdir -p "${CACHE_DIR}"
  release_url="https://dl.k8s.io/${K8S_VERSION}/kubernetes-server-linux-${arch}.tar.gz"
  release_tarball="${CACHE_DIR}/kubernetes-server-${K8S_VERSION}-linux-${arch}.tar.gz"
  expected_checksum="$(curl --fail --location --silent --show-error --retry 5 "${release_url}.sha256")"

  if [[ ! -f "${release_tarball}" ]] || [[ "$(shasum -a 256 "${release_tarball}" | awk '{print $1}')" != "${expected_checksum}" ]]; then
    echo "Downloading the official Kubernetes ${K8S_VERSION} server bundle."
    echo "Interrupted downloads are retried and can be resumed by running this command again."
    curl \
      --fail \
      --location \
      --retry 5 \
      --retry-delay 2 \
      --retry-all-errors \
      --continue-at - \
      --output "${release_tarball}" \
      "${release_url}"
  fi

  actual_checksum="$(shasum -a 256 "${release_tarball}" | awk '{print $1}')"
  if [[ "${actual_checksum}" != "${expected_checksum}" ]]; then
    echo "Checksum verification failed for ${release_tarball}." >&2
    echo "Delete that file and run 'make cluster-up' again." >&2
    exit 1
  fi

  echo "No official Kind node image exists for Kubernetes ${K8S_VERSION}."
  echo "Building ${NODE_IMAGE} from the official Kubernetes release; the first run can take several minutes."
  kind build node-image --type file "${release_tarball}" --image "${NODE_IMAGE}"
else
  echo "Using existing local node image ${NODE_IMAGE}."
fi

echo "Creating Kind cluster '${CLUSTER_NAME}'..."
kind create cluster \
  --name "${CLUSTER_NAME}" \
  --config "${LAB_DIR}/kind-config.yaml" \
  --image "${NODE_IMAGE}" \
  --wait 180s

echo "Limiting every node to ${NODE_CPUS} CPUs and ${NODE_MEMORY} of memory..."
while IFS= read -r node; do
  docker update \
    --cpus "${NODE_CPUS}" \
    --memory "${NODE_MEMORY}" \
    --memory-swap "${NODE_MEMORY}" \
    "${node}" >/dev/null
done < <(kind get nodes --name "${CLUSTER_NAME}")

kubectl config use-context "${CONTEXT}" >/dev/null
"${SCRIPT_DIR}/verify-cluster.sh"
