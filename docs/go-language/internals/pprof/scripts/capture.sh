#!/usr/bin/env bash
set -euo pipefail

profile="${1:-}"
base_url="${2:-http://127.0.0.1:6061}"
seconds="${3:-10}"
output="${4:-}"
root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

usage() {
  echo "usage: $0 {cpu|heap|allocs|goroutine|block|mutex|threadcreate|trace} [pprof-base-url] [seconds] [output]" >&2
  exit 2
}

case "${profile}" in
  cpu)
    endpoint="/debug/pprof/profile?seconds=${seconds}"
    suffix="pb.gz"
    ;;
  heap)
    endpoint="/debug/pprof/heap?gc=1"
    suffix="pb.gz"
    ;;
  allocs)
    endpoint="/debug/pprof/allocs"
    suffix="pb.gz"
    ;;
  goroutine)
    endpoint="/debug/pprof/goroutine"
    suffix="pb.gz"
    ;;
  block)
    endpoint="/debug/pprof/block?seconds=${seconds}"
    suffix="pb.gz"
    ;;
  mutex)
    endpoint="/debug/pprof/mutex?seconds=${seconds}"
    suffix="pb.gz"
    ;;
  threadcreate)
    endpoint="/debug/pprof/threadcreate?seconds=${seconds}"
    suffix="pb.gz"
    ;;
  trace)
    endpoint="/debug/pprof/trace?seconds=${seconds}"
    suffix="trace"
    ;;
  *) usage ;;
esac

mkdir -p "${root_dir}/profiles"
if [[ -z "${output}" ]]; then
  output="${root_dir}/profiles/${profile}.${suffix}"
fi

curl --fail --silent --show-error --output "${output}" "${base_url}${endpoint}"
echo "captured ${profile}: ${output}"
