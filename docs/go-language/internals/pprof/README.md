# Go pprof: a runnable performance-debugging lab

This directory teaches `pprof` by making a service busy, taking real profile
dumps, and following the samples back to source code. It includes:

- `healthy-service/`: the smallest safe HTTP pprof setup and a baseline handler.
- `bad-service/`: bounded, intentional CPU, allocation, heap-retention, I/O,
  goroutine, channel, and mutex problems.
- `loadgen/`: a dependency-free concurrent HTTP load generator.
- `scripts/`: build, capture, and repeatable scenario helpers.
- [ANALYSIS-GUIDE.md](ANALYSIS-GUIDE.md): UI and CLI interpretation, explained
  from first principles.
- [TROUBLESHOOTING-PLAYBOOK.md](TROUBLESHOOTING-PLAYBOOK.md): symptom-to-profile
  decision guide, expected lab findings, false positives, and production advice.

Every bad endpoint has a safety bound, but it still wastes local resources.
Run this only on a development machine.

## 1. The tiny mental model

Imagine stopping the program thousands of times and asking:

> "What function is doing or waiting for the thing I care about right now?"

Each answer is a **sample** containing a call stack and a value such as CPU
time, live bytes, goroutine count, or blocked time. `pprof` groups similar
answers. A large number means "look here," not automatically "this code is
wrong."

```text
symptom -> choose the matching profile -> reproduce representative work
        -> collect samples -> find a wide/hot stack -> inspect your first frame
        -> form a hypothesis -> change code -> repeat under the same workload
```

Do not start with a CPU profile just because it is the famous one. A slow
request waiting for a database uses wall-clock time but almost no CPU.

## 2. Prerequisites and build

Required: Go 1.22 or newer and `curl`. Graph views also need Graphviz (`dot`).

```bash
cd docs/go-language/internals/pprof
go version
command -v curl
command -v dot       # optional for graph rendering
make test
make build
```

The three executables are placed in `bin/`, and generated dumps go in
`profiles/`. Both are ignored by Git.

## 3. How pprof is enabled in Go code

The healthy example uses a separate listener. The important pieces from
[healthy-service/main.go](healthy-service/main.go) are:

```go
import httppprof "net/http/pprof"

runtime.SetBlockProfileRate(1)
runtime.SetMutexProfileFraction(1)

go func() {
    log.Fatal(http.ListenAndServe("127.0.0.1:6060", newPprofMux()))
}()

func newPprofMux() *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /debug/pprof/", httppprof.Index)
    mux.HandleFunc("GET /debug/pprof/cmdline", httppprof.Cmdline)
    mux.HandleFunc("GET /debug/pprof/profile", httppprof.Profile)
    mux.HandleFunc("GET /debug/pprof/symbol", httppprof.Symbol)
    mux.HandleFunc("POST /debug/pprof/symbol", httppprof.Symbol)
    mux.HandleFunc("GET /debug/pprof/trace", httppprof.Trace)
    return mux
}
```

CPU, heap, allocation, goroutine, and thread-creation profiles are available
without enabling a rate. Block and mutex profiles are opt-in. Sampling every
block or mutex event (`1`) is intentionally aggressive for this lab; measure
overhead and use a lower sampling level in production.

### The shortest setup

For a private toy program that does not otherwise use `http.DefaultServeMux`,
this also works:

```go
import (
    "log"
    "net/http"
    _ "net/http/pprof"
)

go func() {
    log.Println(http.ListenAndServe("127.0.0.1:6060", nil))
}()
```

The explicit mux in this lab is preferred because the application port does
not accidentally inherit diagnostic routes. Never expose pprof directly to
the Internet. Profiles can reveal source paths, function names, topology,
traffic shape, and sometimes labels. Bind to loopback, place authentication
and network policy in front of it, or retrieve it through a controlled tunnel.

The official references are [net/http/pprof](https://pkg.go.dev/net/http/pprof)
and [Go diagnostics](https://go.dev/doc/diagnostics).

## 4. First project: profile the healthy service

Use three terminals so load and collection overlap.

### Terminal 1: start the service

```bash
cd docs/go-language/internals/pprof
./scripts/build.sh
./bin/healthy-service
```

Check both listeners:

```bash
curl http://127.0.0.1:8080/healthz
curl http://127.0.0.1:6060/debug/pprof/
```

The app is at `8080`; pprof is at `6060`.

### Terminal 2: generate representative load

```bash
./bin/loadgen \
  -url 'http://127.0.0.1:8080/work?size=50000' \
  -duration 14s \
  -concurrency 8
```

### Terminal 3: capture while the load is active

Start this within the first few seconds of the load:

```bash
./scripts/capture.sh cpu http://127.0.0.1:6060 10 profiles/healthy-cpu.pb.gz
```

This HTTP request stays open for ten seconds while the runtime samples CPU
stacks, then saves a compressed protobuf profile. A profile is not a text log;
keep the binary nearby for the best source and disassembly views.

### Read it in the CLI

```bash
go tool pprof -top ./bin/healthy-service profiles/healthy-cpu.pb.gz
go tool pprof -top -cum ./bin/healthy-service profiles/healthy-cpu.pb.gz
go tool pprof -list='main.buildReport' ./bin/healthy-service profiles/healthy-cpu.pb.gz
```

On the verified run for this tutorial, `main.buildReport` had about 13% flat
CPU and 22% cumulative CPU. Network writes, system calls, and runtime
scheduling were also visible because the local load generator drove about
22,000 requests/second. These exact values will differ by CPU, Go version, and
load; the meaningful baseline is the shape and behavior on your machine.

### Open the pprof web UI

```bash
go tool pprof \
  -http=127.0.0.1:7070 \
  ./bin/healthy-service \
  profiles/healthy-cpu.pb.gz
```

Open `http://127.0.0.1:7070/ui/`. Use the **View** menu to switch between Top,
Graph, Flame Graph, Peek, and Source. Stop the UI with `Ctrl-C`.

If Graph or Flame Graph fails, install Graphviz or use the Top/Source views and
CLI, which are sufficient for the lab.

| View | Question it answers | Use it when |
|---|---|---|
| **Top** | Which functions cost the most? | Always start here |
| **Graph** | How does expensive work flow between callers and callees? | You need ownership or multiple call paths |
| **Flame Graph** | Which complete call stacks dominate? | You want the overall shape of the workload |
| **Peek** | Who calls this function, and what does it call? | Investigating one suspicious function |
| **Source** | Which source line carries the cost? | You have identified your application function |
| **Disassemble** | Which machine instruction carries the cost? | Source-level analysis is insufficient |

#### 1. Top

Top is the best starting point. It ranks functions using a table similar to:

```text
flat     flat%    cum      cum%    function
59.16s   94.69%   59.16s   94.69%  sha256.blockSHA2
0        0%       59.75s   95.63%  main.burnCPU
```

Important columns:

- **Flat**: CPU time spent executing only that function’s own instructions. It excludes functions it calls.
- **Flat%**: is that function’s flat CPU time divided by the profile’s total sampled CPU time.
- **cumulative/cum**: CPU time spent in that function’s own instructions plus all its callees.
```text
A
└── B
    └── C
```
```text
A own instructions = 1s
B own instructions = 2s
C own instructions = 3s
```
| Function | Flat | Cum |
|---|---:|---:|
| `A` | 1s | 6s |
| `B` | 2s | 5s |
| `C` | 3s | 3s |

- **flat%/cum%**: percentage of the entire selected profile.
- **sum%**: running sum of flat percentages down the table.

Eg:
```text
http.serverHandler             ← ancestor, excluded
└── main.workHandler           ← ancestor, excluded
    └── main.buildReport       ← included: 1.12s flat
        ├── fmt.Fprintf        ← included in cum
        └── strings.Builder    ← included in cum
```

Use Top to answer:

- What burns CPU directly?
- What allocates or retains the most memory?
- Where does the most blocking delay accumulate?
- Do three functions explain 90% of the problem, or is cost distributed?

Mental model:

- `flat` = my own bill.
- `cum` = my bill plus everyone I invited.

##### The formula

From the first row, we can reconstruct the total:

```text
syscall.rawsyscalln flat = 4.63 seconds
flat%                  = 48.28%

total CPU samples ≈ 4.63 / 0.4828
                  ≈ 9.59 CPU-seconds
```

Every row uses the same denominator:

```text
Flat% = Flat / 9.59 seconds × 100
```

##### `main.buildReport` in your screenshot

```text
Flat:  1.12s
Flat%: 11.68%
Cum:   1.40s
Cum%:  14.60%
```

This means:

- Approximately `1.12` CPU-seconds were sampled directly in instructions belonging to `buildReport`.
- That is `11.68%` of the profile’s total `9.59` CPU-seconds.
- Another approximately `0.28` CPU-seconds were spent in functions called by `buildReport`.

```text
Cumulative  = Flat + descendants
1.40s       = 1.12s + approximately 0.28s
```

So the rough call shape is:

```text
main.buildReport                 1.40s cumulative
├── its own loop/instructions    1.12s flat
└── called functions             0.28s
```

##### Important caveat

`1.12s` does not mean one call to `buildReport` took 1.12 seconds.

It means approximately 1.12 aggregated CPU-seconds across:

- many requests;
- many calls;
- many goroutines;
- possibly several operating-system threads.

The simplest interpretation is:

> Out of every 100 CPU samples, approximately 12 were taken while the processor was directly executing `main.buildReport`.

Finally, `Sum%` is just a running addition of the sorted `Flat%` column:

```text
48.28%
48.28% + 11.68% = 59.96%
59.96% + 9.38%  = 69.34%
69.34% + 9.07%  = 78.41%  ≈ displayed 78.42%
```

It helps answer: “Your first 4 instruction accounts for x% of cpu time?”

##### Why some rows have zero Flat

For example:

```text
Flat   Flat%   Cum     Cum%    Name
0      0.00%   4.33s   45.15%  syscall.write
```

`syscall.write` appeared in stacks responsible for `4.33` CPU-seconds, but samples were attributed to deeper functions:

```text
syscall.write             flat 0, cumulative 4.33s
└── syscall.syscall       flat 0
    └── syscall.syscalln  flat 0
        └── rawsyscalln   flat 4.63s
```

Therefore:

- `syscall.write` is an owner/ancestor of the work.
- `rawsyscalln` is where the CPU samples landed directly.

Zero flat does **not** mean the function is irrelevant. It may be the application or library entry point responsible for expensive descendants.

##### Why `runtime.kevent` has equal Flat and Cum

```text
Flat: 0.87s
Cum:  0.87s
```

This means all the sampled cost attributed through `runtime.kevent` was also directly attributed to it. There was no additional visible sampled cost in child functions.

#### 2. Graph

Graph displays a call graph:

```text
caller ──sample value──> callee
```
Each node is a function. It normally displays flat and cumulative values. Edges show sampled cost flowing through that call relationship.

- Node properties:
  - Box width represents cumulative share of the selected sample value.
  - value: `flat` of `cum` time
- Edge properties:
  - Thick/red = a large amount of profile value flowed through this call.
  - Thin/grey = a relatively small amount flowed through it.
  - Orange/brown = somewhere between the two.
  - dotted lines = It means pprof removed one or more low-value intermediate nodes while simplifying the graph.

Use Graph when:

- A hot library function has several callers.
- You need to find which endpoint or background job owns a cost.
- One handler calls many expensive operations.
- You need to understand fan-in or fan-out.
- `runtime.mallocgc`, JSON encoding, syscalls, or a lock appears hot and you need its application caller.

Example:

```text
requestHandler
       │
       ├── encodeResponse ──> json.Marshal
       │
       └── makeReport ──────> allocationChurn
```

Graph makes that branching relationship clearer than a flame graph.

Limitation: large programs create crowded “spaghetti graphs.” Focus on a function or package to reduce noise.

#### 3. Flame Graph

This is the view shown in your screenshot.

In this pprof rendering:

- The root is at the top.
- Functions called by it appear underneath.
- Moving downward follows a deeper call stack.
- Box width represents cumulative share of the selected sample value. In the CPU profile if the nodes are wider for a particular call graph -> this means that flow is very expensive and cpu consuming.
- Horizontal position is only stack packing—it is **not time order**.
- Colors mainly distinguish regions/packages. Red does not automatically mean bad.
- The value shown is the its `cumulative` value in that call-stack. (this function + functions called below this box). If the function also has flat samples, pprof displays an additional `self` value.
- Multiple boxes selected when hover? Because the same function can appear in multiple different call-stack paths.

##### Reading your screenshot

The wide path on the left is approximately:

```text
root
└── http.(*conn).serve
    └── http.(*response).finishRequest
        └── bufio.(*Writer).Flush
            └── net.(*conn).Write
                └── internal/poll.(*FD).Write
                    └── syscall.Write
                        └── syscall.rawsyscalln
```

That says a large portion of samples passed through HTTP response writing and operating-system calls.

The application path in the middle is:

```text
root
└── http.(*conn).serve
    └── http.serverHandler.ServeHTTP
        └── http.(*ServeMux).ServeHTTP
            └── http.HandlerFunc.ServeHTTP
                └── main.workHandler
                    └── main.buildReport
```

This connects the generic HTTP machinery to your code. `main.buildReport` is the actionable application function.

The right side contains runtime scheduling work:

```text
runtime.mcall
└── runtime.park_m
    └── runtime.schedule
        └── runtime.findRunnable
```

This is runtime/scheduler activity. It can become visible under high concurrency and does not automatically indicate a scheduler bug.

Use Flame Graph when:

- You want to see the dominant end-to-end call paths.
- You want to understand what sits above a hot leaf.
- Several features call the same expensive operation.
- You want to click and zoom into one wide branch.

The easiest rule is:

> Look for a wide box, then walk upward toward the root to learn who owns it and downward to learn what it calls.

Do not treat the flame graph as an execution timeline. Two boxes beside each other did not necessarily run consecutively.

#### 4. Peek

Peek gives a focused neighborhood around one function:

```text
callers
   ↓
selected function
   ↓
callees
```

Use it when you are asking:

- Who calls `json.Marshal`?
- Which endpoint calls `allocationChurn`?
- What expensive functions does `workHandler` call?
- Is this library function expensive from one caller or several?

Peek sits conceptually between Top and Graph:

- Top shows ranking but little context.
- Graph shows all relationships and can become crowded.
- Peek shows the immediate context around one function.

For your profile, peeking at `main.buildReport` would show how the HTTP handler reaches it and where its cost goes next.

#### 5. Source

Source annotates your actual source lines with flat and cumulative values.

Conceptually:

```go
func buildReport(size int) {
    for i := 0; i < size; i++ {   // 1.2s cumulative
        calculateValue(i)         // 0 flat, 1.1s cumulative
    }
}
```

Use Source after identifying an application function with Top, Graph, Flame Graph, or Peek.

It answers:

- Which loop is expensive?
- Which allocation line retains memory?
- Which channel send is blocking?
- Is time inside this line or inside a function called by it?
- Is a lock held across slow work?

Interpretation:

- High **flat** on a line means execution was sampled directly there.
- Zero flat but high **cumulative** means the function called on that line is expensive.

Source attribution is approximate. Inlining and compiler optimization can place samples on a nearby line.

For your screenshot, open Source for `main.buildReport` after selecting or searching for it.

#### 6. Disassemble

Disassemble maps samples to CPU assembly instructions.

Use it only when you need to answer questions such as:

- Is bounds-checking significant?
- Did the compiler inline or vectorize this function?
- Which atomic or memory instruction is expensive?
- Is a cryptographic/compression routine using an optimized instruction path?
- Why can’t the source view distinguish two operations on the same line?

This is mainly useful for deep CPU investigations. It is rarely useful for goroutine, heap, channel, or ordinary application-level analysis.

### Recommended workflow

For most investigations:

```text
Top
  ↓
Flame Graph
  ↓
Graph or Peek
  ↓
Source
  ↓
Disassemble only if necessary
```

More specifically:

1. Start with **Top** to find the largest values.
2. Open **Flame Graph** to see the complete costly path.
3. Use **Graph** if the hotspot has multiple callers.
4. Use **Peek** to examine one function’s immediate callers and callees.
5. Open **Source** on the first application function you control.
6. Use **Disassemble** only for instruction-level CPU questions.

One final detail: the **SAMPLE** menu changes what the views measure. The same Flame Graph can represent CPU time, live bytes, allocated bytes, object count, goroutine count, or blocking delay. Always check the selected sample type before interpreting box width.

The expanded explanation is also in [ANALYSIS-GUIDE.md](/Users/9999dude/Projects/github/9999dude/omega/docs/go-language/internals/pprof/ANALYSIS-GUIDE.md:41).

## 5. Capture commands worth memorizing

With a service on pprof port `6061`:

```bash
# CPU: sampled on-CPU stacks over an interval
./scripts/capture.sh cpu http://127.0.0.1:6061 30 profiles/cpu.pb.gz

# Current reachable heap, forcing GC immediately before the snapshot
./scripts/capture.sh heap http://127.0.0.1:6061 10 profiles/heap.pb.gz

# Allocation history since process start
./scripts/capture.sh allocs http://127.0.0.1:6061 10 profiles/allocs.pb.gz

# Current goroutine stacks
./scripts/capture.sh goroutine http://127.0.0.1:6061 10 profiles/goroutine.pb.gz

# Synchronization delay during a ten-second interval
./scripts/capture.sh block http://127.0.0.1:6061 10 profiles/block.pb.gz
./scripts/capture.sh mutex http://127.0.0.1:6061 10 profiles/mutex.pb.gz

# Scheduler/runtime event timeline; view with go tool trace, not pprof
./scripts/capture.sh trace http://127.0.0.1:6061 5 profiles/runtime.trace
go tool trace profiles/runtime.trace
```

The HTTP endpoint supports `seconds=N` delta profiles for allocs, block,
goroutine, heap, mutex, and threadcreate, as well as duration-based CPU and
trace capture. `gc=1` asks for GC before a heap snapshot.

For a quick live profile without explicitly saving first:

```bash
go tool pprof 'http://127.0.0.1:6061/debug/pprof/profile?seconds=30'
```

Saving the file explicitly is better for comparisons and incident evidence.

## 6. Second project: deliberately bad code

Start it in Terminal 1:

```bash
./bin/bad-service
```

Application endpoints are on `127.0.0.1:8081`; pprof is on
`127.0.0.1:6061`. In Terminal 2, run one scenario:

```bash
./scripts/scenario.sh cpu
./scripts/scenario.sh allocs
./scripts/scenario.sh heap
./scripts/scenario.sh io
./scripts/scenario.sh goroutine
./scripts/scenario.sh channel
./scripts/scenario.sh mutex
```

Run **one scenario at a time**. Restart `bad-service` between scenarios when
you want a clean experiment: heap and allocation counters are process-wide,
and leaked goroutines intentionally remain until restart. The script prints
request rate and latency and writes the corresponding dump under `profiles/`.

Inspect any binary profile in the same UI:

```bash
go tool pprof -http=127.0.0.1:7070 ./bin/bad-service profiles/cpu.pb.gz
```

Profile-specific commands and the observed findings are in
[TROUBLESHOOTING-PLAYBOOK.md](TROUBLESHOOTING-PLAYBOOK.md).

## 7. A disciplined investigation loop

1. Write the symptom numerically: "p99 rose from 80 ms to 900 ms while CPU
   stayed at 25%," not merely "the app is slow."
2. Check external bounds first: CPU quota/throttling, memory limit/OOM, disk or
   network saturation, downstream latency, request mix, and recent deploys.
3. Pick the profile whose sample value matches the symptom.
4. Warm the process and reproduce realistic traffic. Idle profiles answer
   questions about idle code.
5. Collect only one expensive diagnostic at a time and record duration, build,
   instance, traffic, and runtime settings.
6. In `top`, find a dominant leaf by **flat** value. In `top -cum`, find the
   request/job path responsible by **cumulative** value.
7. Walk from runtime/library frames toward the first frame owned by your code.
8. Use `list` or Source to find the exact line. Read the surrounding ownership,
   cancellation, and lifetime logic; do not optimize a line in isolation.
9. State a falsifiable hypothesis, change one thing, rerun the same workload,
   and compare throughput, latency, resource metrics, and profiles.
10. Keep the fix only if user-visible behavior improves without moving the
    bottleneck or breaking correctness.

`pprof` tells you **where cost accumulated**. It does not tell you what the
correct architecture is, whether the result is valuable, or whether a lock,
allocation, or goroutine is logically necessary.

## 8. Exercises

1. Predict the widest CPU leaf before running `scenario.sh cpu`; confirm with
   `top` and `list`.
2. Run allocs, compare `-alloc_space` with `-inuse_space`, and explain how an
   app can have severe GC pressure without a large live heap.
3. Capture heap before and after 32 `/leak` calls, then use `-diff_base`.
4. Run I/O load and explain why the CPU profile can look quiet while requests
   take about one second.
5. Find the repeated `[chan receive]` stack in the plaintext goroutine dump.
6. Explain why `runtime.chansend1` is a symptom but `sendWithBackpressure` is
   the first actionable application frame.
7. Explain why a mutex profile reports `Unlock` and the holder stack rather
   than simply blaming `Lock`.
8. Fix one problem, preserve the load parameters, and prove improvement using
   both service-level latency and a new profile.

## 9. Cleanup

Stop services and UI processes with `Ctrl-C`. Generated artifacts can be
removed with:

```bash
make clean
```
