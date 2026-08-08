#!/usr/bin/env bash
set -euo pipefail

scenario="${1:-}"
root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
app_url="${APP_URL:-http://127.0.0.1:8081}"
pprof_url="${PPROF_URL:-http://127.0.0.1:6061}"
load="${root_dir}/bin/loadgen"
capture="${root_dir}/scripts/capture.sh"

if [[ ! -x "${load}" ]]; then
  "${root_dir}/scripts/build.sh"
fi

case "${scenario}" in
  cpu)
    "${load}" -url "${app_url}/cpu?ms=100" -duration 14s -concurrency 8 &
    load_pid=$!
    sleep 1
    "${capture}" cpu "${pprof_url}" 10 "${root_dir}/profiles/cpu.pb.gz"
    wait "${load_pid}"
    ;;
  allocs)
    "${load}" -url "${app_url}/alloc?kb=64&rounds=32" -duration 8s -concurrency 8
    "${capture}" allocs "${pprof_url}" 10 "${root_dir}/profiles/allocs.pb.gz"
    ;;
  heap)
    for _ in {1..32}; do
      curl --fail --silent --show-error "${app_url}/leak?kb=512" >/dev/null
    done
    "${capture}" heap "${pprof_url}" 10 "${root_dir}/profiles/heap.pb.gz"
    ;;
  io)
    "${load}" -url "${app_url}/io?ms=1000" -duration 8s -concurrency 32 &
    load_pid=$!
    sleep 2
    "${capture}" goroutine "${pprof_url}" 10 "${root_dir}/profiles/io-goroutine.pb.gz"
    "${capture}" trace "${pprof_url}" 3 "${root_dir}/profiles/io.trace"
    wait "${load_pid}"
    ;;
  goroutine)
    for _ in {1..20}; do
      curl --fail --silent --show-error "${app_url}/goroutine-leak?n=25" >/dev/null
    done
    "${capture}" goroutine "${pprof_url}" 10 "${root_dir}/profiles/goroutine-leak.pb.gz"
    ;;
  channel)
    "${load}" -url "${app_url}/channel?messages=20&consumer_ms=10" -duration 14s -concurrency 16 &
    load_pid=$!
    sleep 1
    "${capture}" block "${pprof_url}" 10 "${root_dir}/profiles/block.pb.gz"
    wait "${load_pid}"
    ;;
  mutex)
    "${load}" -url "${app_url}/mutex?workers=16&loops=10&hold_ms=1" -duration 14s -concurrency 8 &
    load_pid=$!
    sleep 1
    "${capture}" mutex "${pprof_url}" 10 "${root_dir}/profiles/mutex.pb.gz"
    wait "${load_pid}"
    ;;
  *)
    echo "usage: $0 {cpu|allocs|heap|io|goroutine|channel|mutex}" >&2
    exit 2
    ;;
esac
