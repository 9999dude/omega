# Troubleshooting playbook and bad-service walkthrough

The numbers below came from a verified run on Go 1.26.5 on macOS/arm64. They
are clues to reproduce, not golden test values.

## 1. Which profile do I take?

| Observed symptom | First evidence | First profile | Follow-up |
|---|---|---|---|
| CPU cores saturated, throttling, throughput falls | CPU usage/quota and request rate | CPU, 20-30 seconds | `top`, `top -cum`, Source, OS profiler for cgo/kernel |
| RSS or Go heap grows after traffic | post-GC heap/runtime metrics | heap `?gc=1`, repeated snapshots | in-use space/objects, diff, GC metrics |
| GC CPU high but live heap modest | allocation rate and GC cycles | allocs delta / alloc-space | alloc objects, CPU profile |
| High latency with low CPU | latency plus dependency/OS metrics | goroutine snapshot during incident | execution trace, distributed trace, downstream metrics |
| Goroutine count grows and does not fall | runtime goroutine metric | goroutine, twice over time | plaintext `debug=2`, ownership/cancellation review |
| Producers/consumers stall | queue depth and latency | block profile | goroutine stacks, execution trace |
| Low CPU utilization under concurrency | contention/throughput evidence | mutex and block profiles | critical-section source, trace |
| OS thread count explodes | process/thread metrics | threadcreate | goroutine stacks, cgo and `LockOSThread` review |

Do not choose from the symptom name alone. "Memory issue" could mean retention,
allocation churn, goroutine stacks, cgo memory, page cache, or an intentionally
large cache.

## 2. CPU scenario

Run against a freshly started bad service:

```bash
./scripts/scenario.sh cpu
go tool pprof -top ./bin/bad-service profiles/cpu.pb.gz
go tool pprof -top -cum ./bin/bad-service profiles/cpu.pb.gz
go tool pprof -list='main.burnCPU' ./bin/bad-service profiles/cpu.pb.gz
```

Verified signature:

```text
flat      flat%   cum       function
59.16s    94.69%  59.16s    crypto/.../sha256.blockSHA2
0          0%     59.75s    main.burnCPU
```

How to reason:

1. High flat value says the SHA-256 block implementation is where CPU executes.
2. `top -cum` says `main.burnCPU` owns about 96% of the sampled call path.
3. Source says the call to `sha256.Sum256(data)` carries almost all cumulative
   cost.
4. The runtime/library code is behaving normally; the application bug is
   repeatedly hashing 32 KiB until a deadline for no useful reason.

Watch for:

- a hot leaf called by several features—use Peek/Graph to find the important
  caller;
- CPU quota: one fully used core can display as 100%, 25%, or throttled
  depending on dashboard denominator;
- short captures with few samples—collect longer under stable load;
- compiler inlining and runtime frames—follow the call path, not just names;
- high cumulative router/worker functions with zero flat value—ownership, not
  leaf cost;
- CPU profile missing waits—CPU intentionally excludes sleep and I/O wait.

Possible fixes depend on real intent: remove duplicate work, cache bounded
results, use a better algorithm/data structure, batch, reduce encode/copy
passes, or move optional work off the critical path. Reprofile before adding
parallelism; parallel waste can worsen CPU and tail latency.

## 3. Allocation churn versus heap retention

### Churn

```bash
./scripts/scenario.sh allocs
go tool pprof -top -alloc_space ./bin/bad-service profiles/allocs.pb.gz
go tool pprof -top -inuse_space ./bin/bad-service profiles/allocs.pb.gz
```

Verified signature after eight seconds:

```text
alloc_space: main.allocationChurn   about 47,793 MB (99.58%)
inuse_space: total profile          about      4 MB
```

This is not a 48 GB leak. It means about 48 GB flowed through sampled heap
allocations over the process lifetime while GC reclaimed nearly all of it. The
symptoms in a real service would be allocation rate, GC CPU, memory bandwidth,
and latency—not necessarily a steadily rising post-GC heap.

Watch `alloc_space` for:

- `make([]T, ...)` inside loops;
- string/`[]byte` conversions and concatenation;
- reflection, interface boxing, JSON/protobuf intermediates;
- request buffering and repeated decompression;
- per-item objects that could be one contiguous structure;
- pools that reduce allocations but retain too much memory or add contention.

Prefer ownership and algorithm changes before `sync.Pool`. Validate with
benchmarks (`-benchmem`), alloc-space, GC metrics, and service latency.

### Retention/leak

Restart the service first, then:

```bash
./scripts/scenario.sh heap
go tool pprof -top -inuse_space ./bin/bad-service profiles/heap.pb.gz
go tool pprof -list='main.*leak' -inuse_space \
  ./bin/bad-service profiles/heap.pb.gz
curl http://127.0.0.1:8081/stats
```

The scenario retains 32 x 512 KiB buffers. Verified post-GC signature:

```text
main.(*problemApp).leak   about 17.8 MB, 81% of sampled live heap
stats.retained_bytes      16,777,216 bytes
```

The sampled number is approximate and includes allocator/profile behavior; the
explicit application counter is exact. The important evidence is that the
allocation path remains reachable after GC and grows with repeated calls.

Common retention roots:

- unbounded maps/caches/dedup tables;
- slices that keep a huge backing array for a tiny subslice;
- queues/channels whose consumers lag;
- timers/tickers or callbacks never stopped;
- goroutines retaining request objects in their stacks;
- global registries and metrics labels with unbounded cardinality;
- response/request bodies or buffers held past request completion.

Take at least two post-GC snapshots across equal work. A single large heap can
be a legitimate steady-state cache. Growth and ownership are the evidence.

If RSS grows but Go in-use heap does not, investigate goroutine stacks, runtime
metadata, fragmentation/scavenging, mmap/page cache, cgo/native libraries, and
the container/OS measurement definition.

## 4. I/O and wall-latency scenario

Restart to remove prior leaked goroutines, then:

```bash
./scripts/scenario.sh io
go tool pprof -top ./bin/bad-service profiles/io-goroutine.pb.gz
go tool trace profiles/io.trace
```

The handler waits one second for a deliberately slow pipe. Under concurrency
32, the load generator reports roughly 32 requests/second and one-second
average latency. A goroutine snapshot taken during load shows about 32 stacks
through:

```text
internal/poll.(*FD).Read
io.ReadAll
main.readSlowDependency
main.(*problemApp).slowIO
```

and about 32 paired `main.delayedWrite` sleepers.

Why CPU may look innocent: the waiting goroutines are parked. They consume
wall-clock latency and retain stack/request resources but do not continuously
execute on CPU.

For real network/disk I/O, correlate:

- goroutine state (`IO wait`, `syscall`, `semacquire`, channel states);
- client/server timeouts and context deadlines;
- DNS, connect, TLS, time-to-first-byte, body read, and pool wait timings;
- connection-pool limits and file descriptor exhaustion;
- disk latency/queue depth or network retransmits/saturation;
- downstream latency/error/queue metrics;
- execution trace and distributed spans.

The Go block profile focuses on synchronization primitives; do not expect it
to explain every network or disk wait. Goroutine stacks show where waiters are
now, while trace shows when scheduler, syscall, GC, and blocking events happen.

## 5. Goroutine leak scenario

Restart, then:

```bash
./scripts/scenario.sh goroutine
go tool pprof -top ./bin/bad-service profiles/goroutine-leak.pb.gz
curl 'http://127.0.0.1:6061/debug/pprof/goroutine?debug=2'
```

Verified signature:

```text
500 goroutines (about 99% of stacks) -> main.leakedWorker
state: [chan receive]
source: <-neverReleased
```

Binary output is good for grouping and graphing. `debug=2` plaintext is often
best for reading exact goroutine states and creation stacks.

Mental model:

- one blocked goroutine is normal;
- many identical blocked goroutines may be a pool or a leak;
- a count that grows with traffic and does not return after work completes is
  strong leak evidence;
- the blocking line is only half the story—find who created it, who owns its
  stop condition, and why cancellation/close/send never occurs.

Audit every long-lived goroutine for:

- owner and lifetime;
- `context.Context` cancellation;
- channel close/send responsibility;
- timeout and retry limits;
- error/early-return paths;
- ticker/timer cleanup;
- `WaitGroup` completion;
- bounded fan-out and queues.

Do not close a channel merely to make the profile disappear. Channel closing is
an ownership protocol, and the wrong closer can panic or lose work.

## 6. Channel/backpressure scenario

Restart, then:

```bash
./scripts/scenario.sh channel
go tool pprof -top -sample_index=delay \
  ./bin/bad-service profiles/block.pb.gz
go tool pprof -top -cum -sample_index=delay \
  ./bin/bad-service profiles/block.pb.gz
```

Verified signature for a ten-second delta profile:

```text
runtime.chansend1          about 151s delay (89%)
main.sendWithBackpressure  about 159s cumulative delay (94%)
```

Delay exceeds ten seconds because many handlers wait at once. The runtime send
is the mechanism; `main.sendWithBackpressure` contains the design decision: an
unbuffered queue feeding a slow consumer.

Questions to ask:

- Is the consumer slower, dead, or waiting on its own dependency?
- Is blocking intentional backpressure or an accidental throughput ceiling?
- Would a small bounded buffer absorb legitimate bursts?
- What happens when the buffer fills: block, reject, shed, spill, or drop?
- Does work hold large objects while queued?
- Is there a missing receiver, send, close, or cancellation path?
- Is `select { default: }` silently dropping important work?

A huge channel buffer hides overload, increases memory retention, and delays
failure. Fix service rate and overload policy; size a bounded queue from the
allowed burst and memory budget.

## 7. Mutex contention scenario

Restart, then:

```bash
./scripts/scenario.sh mutex
go tool pprof -top -sample_index=delay \
  ./bin/bad-service profiles/mutex.pb.gz
go tool pprof -list='main.*holdHotLock' -sample_index=delay \
  ./bin/bad-service profiles/mutex.pb.gz
```

Verified signature:

```text
sync.(*Mutex).Unlock             about 1,255s delay (100%)
main.(*problemApp).holdHotLock   about 1,255s cumulative delay
```

The mutex profile attributes contention to the stack that released the lock,
because that holder made other goroutines wait. Seeing `Unlock` is expected.
Source reveals the real bug: `time.Sleep` runs while the global mutex is held.

Questions to ask:

- Can slow I/O, logging, encoding, callbacks, or sleep move outside the lock?
- Is one global lock protecting independent keys that could be sharded?
- Can immutable snapshots or ownership-by-one-goroutine simplify sharing?
- Is an `RWMutex` genuinely read-heavy with short reads, or will writer
  starvation/overhead make it worse?
- Is contention significant under representative concurrency, or only in a
  stress test far beyond capacity?
- Does removing the lock introduce races? Always run correctness and race tests.

Never "fix" a mutex profile by deleting synchronization without redesigning
ownership. A fast data race is still a bug.

## 8. Thread and scheduler issues

Useful commands:

```bash
./scripts/capture.sh goroutine http://127.0.0.1:6061 10 profiles/goroutine.pb.gz
./scripts/capture.sh threadcreate http://127.0.0.1:6061 10 profiles/threadcreate.pb.gz
./scripts/capture.sh trace http://127.0.0.1:6061 5 profiles/scheduler.trace
go tool trace profiles/scheduler.trace
```

Goroutines are not OS threads. Large goroutine count does not necessarily mean
large thread count. Threads can grow around blocking syscalls, cgo,
`runtime.LockOSThread`, or scheduler needs. Correlate threadcreate with process
thread metrics and goroutine states.

Execution trace is the better microscope for:

- runnable goroutines not getting CPU;
- long stop-the-world or GC-assist effects;
- serialization and scheduler latency;
- syscall blocking and unblocking;
- goroutine creation/unblocking chains;
- short time-correlated latency events.

Trace is detailed and higher-overhead than ordinary profiles. Capture short,
focused windows and avoid collecting other expensive diagnostics at the same
time.

## 9. Production checklist

Before capture:

- identify build/version, instance, time window, symptom, and traffic mix;
- confirm CPU/memory quota and external saturation;
- use a representative instance and warm process;
- secure diagnostic access;
- estimate overhead and choose rates/duration;
- collect one profile type at a time.

During analysis:

- confirm profile type and sample index at the top of output;
- check duration and total sample value—is there enough signal?
- use flat for expensive leaves and cumulative for ownership;
- find the first actionable application frame;
- compare repeated snapshots for leaks;
- treat runtime/library frames as mechanisms until call paths prove otherwise;
- keep a hypothesis/evidence log instead of random code changes.

After a change:

- rerun equal work with the same concurrency and environment;
- compare throughput, p50/p95/p99, errors, CPU, heap, allocation rate, GC, and
  goroutine count;
- inspect the new profile for a shifted bottleneck;
- run tests and the race detector where concurrency changed;
- retain profile/build metadata needed to reproduce the result.

## 10. Frequent profiling mistakes

- Profiling idle or synthetic-unrepresentative traffic.
- Treating a single snapshot as proof of a leak.
- Reading alloc-space as current memory.
- Reading a flame graph as a timeline.
- Assuming high cumulative cost means high self cost.
- Blaming runtime functions without walking to their application callers.
- Expecting CPU profiles to contain sleep and I/O wait.
- Expecting block profiles to explain all I/O.
- Ignoring container throttling, downstreams, or kernel/device limits.
- Comparing different builds, traffic, warmup, durations, or amounts of work.
- Enabling maximum mutex/block sampling everywhere without measuring overhead.
- Publishing pprof endpoints or incident dumps without access control.
- Optimizing profile percentages without verifying user-visible improvement.

The goal is not a pretty profile. The goal is a tested explanation connecting
symptom -> resource -> stack -> source decision -> measurable fix.
