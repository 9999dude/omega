#!/usr/bin/env bash

set -euo pipefail

KIND_VERSION="v0.32.0"
KUBECTL_VERSION="v1.36.2"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LAB_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
BIN_DIR="${LAB_DIR}/.tools/bin"

require() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    exit 1
  fi
}

download_and_verify() {
  local url="$1"
  local checksum_url="$2"
  local target="$3"
  local temp_file
  local expected_checksum
  local actual_checksum

  temp_file="$(mktemp "${target}.XXXXXX")"
  trap 'rm -f "${temp_file}"' RETURN
  curl --fail --location --silent --show-error "${url}" --output "${temp_file}"
  expected_checksum="$(curl --fail --location --silent --show-error "${checksum_url}" | awk '{print $1}')"
  actual_checksum="$(shasum -a 256 "${temp_file}" | awk '{print $1}')"

  if [[ "${expected_checksum}" != "${actual_checksum}" ]]; then
    rm -f "${temp_file}"
    echo "Checksum verification failed for ${url}" >&2
    exit 1
  fi

  chmod +x "${temp_file}"
  mv "${temp_file}" "${target}"
  trap - RETURN
}

require curl
require shasum

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "${os}" in
  darwin|linux) ;;
  *)
    echo "Unsupported operating system: ${os}" >&2
    exit 1
    ;;
esac

machine="$(uname -m)"
case "${machine}" in
  x86_64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *)
    echo "Unsupported CPU architecture: ${machine}" >&2
    exit 1
    ;;
esac

mkdir -p "${BIN_DIR}"

kind_binary="${BIN_DIR}/kind"
if [[ ! -x "${kind_binary}" ]] || ! "${kind_binary}" version 2>/dev/null | grep -q "${KIND_VERSION}"; then
  echo "Installing Kind ${KIND_VERSION}..."
  kind_url="https://github.com/kubernetes-sigs/kind/releases/download/${KIND_VERSION}/kind-${os}-${arch}"
  download_and_verify "${kind_url}" "${kind_url}.sha256sum" "${kind_binary}"
fi

kubectl_binary="${BIN_DIR}/kubectl"
if [[ ! -x "${kubectl_binary}" ]] || ! "${kubectl_binary}" version --client=true --output=yaml 2>/dev/null | grep -q "gitVersion: ${KUBECTL_VERSION}"; then
  echo "Installing kubectl ${KUBECTL_VERSION}..."
  kubectl_url="https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/${os}/${arch}/kubectl"
  download_and_verify "${kubectl_url}" "${kubectl_url}.sha256" "${kubectl_binary}"
fi

echo "Local tools are ready in ${BIN_DIR}"
