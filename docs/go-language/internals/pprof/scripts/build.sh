#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
mkdir -p "${root_dir}/bin"

go build -o "${root_dir}/bin/healthy-service" "${root_dir}/healthy-service"
go build -o "${root_dir}/bin/bad-service" "${root_dir}/bad-service"
go build -o "${root_dir}/bin/loadgen" "${root_dir}/loadgen"

echo "built binaries in ${root_dir}/bin"
