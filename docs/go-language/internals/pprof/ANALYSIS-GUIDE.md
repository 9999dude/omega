# How to read pprof without guessing

This guide explains the numbers and pictures. Work through it with one of the
generated profiles open.

## 1. A profile is a weighted collection of stack traces

Suppose the profiler records these simplified CPU samples:

```text
3 samples: HTTP -> handler -> parse
7 samples: HTTP -> handler -> hash
```

Then:

- `parse` has flat value 3 and cumulative value 3.
- `hash` has flat value 7 and cumulative value 7.
- `handler` has flat value 0 and cumulative value 10.
- `HTTP` has flat value 0 and cumulative value 10.

The parents did not burn CPU themselves, but they caused children to do it.
This is the single most important pprof idea.

### Flat means "in this function itself"

High flat cost asks: **what leaf operation is expensive?**

Examples: hashing, encoding, copying, allocating, a runtime channel send, or a
mutex unlock carrying contention samples.

### Cumulative means "this function plus everything it called"

High cumulative cost asks: **which higher-level feature/request/job owns the
expensive operation?**

A router or generic worker loop often has huge cumulative cost but almost zero
flat cost. Removing the router is not the conclusion; walk down its costly
children.

## 2. Reading the Top table like a sentence

Example shape:

```text
      flat  flat%   sum%        cum   cum%
    59.16s 94.69% 94.69%     59.16s 94.69%  crypto/...blockSHA2
         0     0% 94.69%     59.75s 95.63%  main.burnCPU
```

- `flat`: sampled value directly in this function.
- `flat%`: flat divided by total sampled value.
- `sum%`: running total of `flat%` down the currently sorted table. It helps
  answer "how many top rows explain 80%?"
- `cum`: this function plus descendants.
- `cum%`: cumulative value divided by total.
- `Total samples`: denominator for the percentages.
- `Dropped N nodes`: pprof hid low-value nodes using its display threshold; it
  did not lose them from the profile. Increase node count or change focus when
  the hidden tail matters.

Run both sorts:

```bash
go tool pprof -top ./bin/bad-service profiles/cpu.pb.gz
go tool pprof -top -cum ./bin/bad-service profiles/cpu.pb.gz
```

Baby version:

- `flat` = "my own bill."
- `cum` = "my bill plus everybody I invited."

For CPU profiles, sampled CPU time is summed across executing threads. On a
multicore machine, total samples can exceed elapsed wall time. A 10-second
capture showing 60 CPU-seconds means the process averaged roughly six logical
CPUs during the capture, not that time travel occurred.

For block and mutex profiles, summed delay can greatly exceed wall time because
many goroutines wait concurrently. Sixteen goroutines each waiting one second
produce about sixteen goroutine-seconds of delay.

## 3. Reading the Graph view

A graph node is a function. An arrow means caller -> callee. The arrow label or
width represents sampled value flowing along that call path.

Each node normally shows:

```text
function name
flat value (flat percent)
cumulative value (cumulative percent)
```

Use visual size/color to find candidates, but use the printed values as truth;
exact visual mapping varies with the selected view and positive/negative
differential profiles.

How to walk a graph:

1. Find a large leaf with high flat value.
2. Follow arrows upward until reaching the first meaningful application frame.
3. Follow arrows downward from a high-cumulative application frame to see what
   it delegated to.
4. Click a node to focus. Use **Refine -> Reset** when filtering makes the
   picture mysteriously empty.

Common mistake: blaming `runtime.mallocgc`, `runtime.chansend1`, `syscall`, or
`encoding/json` merely because it is large. That runtime/library frame names
the mechanism. Its application caller often identifies the decision that
created the cost.

## 4. Reading a Flame Graph

A flame graph is the same stacks arranged differently:

- vertical direction: call depth; callers below, callees above;
- horizontal width: share of the selected sample value;
- horizontal position: packing, not chronological order;
- a wide box: this function appears in many costly stacks cumulatively;
- a wide top/leaf box: much of the value is directly in that function;
- color: mainly visual separation; width carries the important quantity.

Click a wide application frame to zoom into one ancestry. Search for a function
name when the graph is crowded. A flame graph is not a trace timeline; adjacent
boxes did not necessarily execute next to each other in time.

## 5. Peek and Source

**Peek** shows callers and callees around a selected function. Use it when a
generic library leaf has several callers and you need the expensive one.

**Source** (CLI: `list`) annotates source lines with flat and cumulative values:

```bash
go tool pprof \
  -list='main.burnCPU' \
  ./bin/bad-service \
  profiles/cpu.pb.gz
```

Interpret source values carefully:

- cost is assigned to sampled program counters and can land on a nearby line;
- inlining and compiler optimization can move or merge apparent work;
- a call line with zero flat but huge cumulative value means the callee is
  expensive;
- missing source usually means the local binary/source does not match the
  profile, paths moved, or symbols were stripped.

For instruction-level problems:

```bash
go tool pprof -disasm='regexp' ./bin/bad-service profiles/cpu.pb.gz
```

Do this after source-level analysis, not as the first step.

## 6. Interactive CLI cheat sheet

Start an interactive session:

```bash
go tool pprof ./bin/bad-service profiles/cpu.pb.gz
```

At the `(pprof)` prompt:

```text
top                 hottest flat functions
top -cum            hottest cumulative call paths
top 20              show more rows
list burnCPU        annotated source matching regexp
peek burnCPU        callers and callees
web                 open a Graphviz graph
svg                 write an SVG
tags                list profile labels/tags, if any
help                all commands and options
quit                exit
```

Useful non-interactive filtering:

```bash
# Keep stacks containing application functions.
go tool pprof -top -focus='main\.' ./bin/bad-service profiles/cpu.pb.gz

# Remove noisy display frames only after first viewing the full profile.
go tool pprof -top -ignore='runtime\.' ./bin/bad-service profiles/cpu.pb.gz

# Increase how many nodes survive graph trimming.
go tool pprof -http=127.0.0.1:7070 -nodecount=200 \
  ./bin/bad-service profiles/cpu.pb.gz
```

`focus` keeps complete stacks that match; it does not magically reattribute
runtime cost to your code. Filtering too early can hide the actual bottleneck.

## 7. Choose the sample type before interpreting memory

One memory dump can carry four useful sample indexes:

- `inuse_space`: live sampled bytes (default for heap questions).
- `inuse_objects`: live sampled object count.
- `alloc_space`: sampled bytes allocated over the interval/process lifetime.
- `alloc_objects`: sampled allocation count over the interval/process lifetime.

Commands:

```bash
go tool pprof -top -inuse_space  ./bin/bad-service profiles/heap.pb.gz
go tool pprof -top -inuse_objects ./bin/bad-service profiles/heap.pb.gz
go tool pprof -top -alloc_space  ./bin/bad-service profiles/allocs.pb.gz
go tool pprof -top -alloc_objects ./bin/bad-service profiles/allocs.pb.gz
```

Baby analogy:

- in-use = toys still on the floor now;
- alloc = every toy brought through the room, even if put away.

A retention leak is large and growing in-use space after GC. Allocation churn
can have enormous alloc space, small in-use space, frequent GC, and high CPU.
They are different bugs with different fixes.

Memory profiles are sampled, so small objects and exact byte counts vary. Heap
is Go heap, not total RSS: goroutine stacks, runtime metadata, code mappings,
native/cgo allocations, and OS behavior can make process RSS disagree.

## 8. Profile type mental models

| Profile | Each sample roughly means | Good first question | Does not prove |
|---|---|---|---|
| CPU | sampled time while executing on CPU | Which leaves consume cycles? | Why wall latency is high while waiting |
| heap / in-use | reachable sampled heap at snapshot | Who retains live bytes/objects? | That high allocation rate exists |
| allocs / alloc | sampled allocation history | Who creates GC work? | That those objects remain live |
| goroutine | one current goroutine stack | Which stacks repeat, and are counts growing? | How much CPU each goroutine uses |
| block | cumulative synchronization waiting | Where do sends/receives/select/waits block? | General network/disk I/O latency |
| mutex | cumulative contention attributed to holders | Which critical section makes others wait? | That every lock is bad |
| threadcreate | stack that caused OS-thread creation | Why are many OS threads created? | Total goroutine count |
| trace | timestamped scheduler/GC/block/syscall events | When and why was work delayed? | Long-term statistically representative cost |

## 9. Comparing two profiles

Absolute profiles answer "where did cost go?" Differential profiles answer
"what changed?"

```bash
# Capture a warmed-up baseline.
./scripts/capture.sh heap http://127.0.0.1:6061 10 profiles/heap-before.pb.gz

# Perform a fixed amount of work, then capture again.
./scripts/capture.sh heap http://127.0.0.1:6061 10 profiles/heap-after.pb.gz

go tool pprof -top \
  -diff_base=profiles/heap-before.pb.gz \
  ./bin/bad-service profiles/heap-after.pb.gz
```

Positive rows grew; negative rows shrank relative to the baseline. A diff is
only fair when binaries, traffic mix, warmup, duration, sample settings, and
amount of work are comparable. For performance comparisons, equal work is
often more meaningful than equal wall duration.

## 10. Spotting bugs quickly

Use these visual smells as leads:

- one CPU leaf owns most samples: an algorithmic hotspot, tight loop, encoding,
  copying, hashing, compression, or lock spinning;
- flat CPU scattered through allocator/GC frames plus huge alloc-space: churn;
- one application allocation path grows across post-GC heap snapshots:
  retention/leak/cache/queue growth;
- thousands of identical goroutine stacks: leak, stuck fan-out, missing timeout,
  or intentional worker pool—validate lifecycle;
- many stacks at channel send: consumer slower than producer or missing
  backpressure policy;
- many stacks at channel receive: missing producer/close, idle workers, or leak;
- block delay concentrated in `select` or channel operations: synchronization
  is gating throughput;
- mutex delay at one unlock stack: the holder's critical section is too broad
  or overly shared;
- low CPU plus many `IO wait` stacks and high latency: downstream/disk/network
  wait; use timeouts, metrics, tracing, and execution trace;
- runtime/syscall leaves dominate: follow upward to application ownership and
  check OS/container metrics before rewriting runtime code.

Always ask whether the cost is proportional to useful throughput. A function
using 60% CPU while doing 90% of valuable work may be healthy. The same 60% in
duplicate serialization may be a bug.

## 11. What pprof cannot tell you alone

- which user/request/tenant was harmed unless you attach labels;
- end-to-end latency across services;
- queue time outside the process;
- kernel, device, network, or downstream saturation in full detail;
- correctness of goroutine/channel ownership;
- whether a retained cache is intentional and properly bounded;
- whether an optimization improves p95/p99 latency or only a microbenchmark.

Correlate profiles with request rate, error rate, latency histograms, CPU quota
and throttling, GC/runtime metrics, RSS, file descriptors, disk/network metrics,
downstream telemetry, logs, and distributed traces.
