package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type problemApp struct {
	maxRetainedBytes int

	heapMu       sync.Mutex
	retainedHeap [][]byte
	retained     int

	neverReleased  chan struct{}
	leakedRoutines atomic.Int64

	hotMu         sync.Mutex
	sharedCounter int64
}

func newProblemApp(maxRetainedBytes int) *problemApp {
	return &problemApp{
		maxRetainedBytes: maxRetainedBytes,
		neverReleased:    make(chan struct{}),
	}
}

func (a *problemApp) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", a.index)
	mux.HandleFunc("GET /healthz", textHandler("ok\n"))
	mux.HandleFunc("GET /cpu", a.cpu)
	mux.HandleFunc("GET /alloc", a.alloc)
	mux.HandleFunc("GET /leak", a.leak)
	mux.HandleFunc("GET /io", a.slowIO)
	mux.HandleFunc("GET /goroutine-leak", a.goroutineLeak)
	mux.HandleFunc("GET /channel", a.channelContention)
	mux.HandleFunc("GET /mutex", a.mutexContention)
	mux.HandleFunc("GET /stats", a.stats)
	return mux
}

func (a *problemApp) index(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = io.WriteString(w, `Intentional pprof lab problems:
/cpu?ms=100
/alloc?kb=64&rounds=32
/leak?kb=512
/io?ms=500
/goroutine-leak?n=25
/channel?messages=10&consumer_ms=25
/mutex?workers=16&loops=20&hold_ms=2
/stats
`)
}

func (a *problemApp) cpu(w http.ResponseWriter, r *http.Request) {
	duration := time.Duration(queryInt(r, "ms", 100, 1, 2_000)) * time.Millisecond
	iterations, digest := burnCPU(duration)
	writeJSON(w, map[string]any{"iterations": iterations, "digest_prefix": fmt.Sprintf("%x", digest[:4])})
}

//go:noinline
func burnCPU(duration time.Duration) (int, [sha256.Size]byte) {
	deadline := time.Now().Add(duration)
	data := make([]byte, 32*1024)
	var digest [sha256.Size]byte
	iterations := 0
	// BUG: expensive work scales with wall time and has no useful result.
	for time.Now().Before(deadline) {
		digest = sha256.Sum256(data)
		data[iterations%len(data)] ^= digest[0]
		iterations++
	}
	return iterations, digest
}

func (a *problemApp) alloc(w http.ResponseWriter, r *http.Request) {
	kb := queryInt(r, "kb", 64, 1, 1_024)
	rounds := queryInt(r, "rounds", 32, 1, 1_000)
	bytes := allocationChurn(kb<<10, rounds)
	writeJSON(w, map[string]any{"temporary_bytes": bytes})
}

//go:noinline
func allocationChurn(bytesPerRound, rounds int) int {
	total := 0
	for i := 0; i < rounds; i++ {
		// BUG: high-rate temporary allocation adds GC pressure even though the
		// live heap can remain small.
		buffer := make([]byte, bytesPerRound)
		for page := 0; page < len(buffer); page += 4096 {
			buffer[page] = byte(i)
		}
		total += len(buffer)
		runtime.KeepAlive(buffer)
	}
	return total
}

func (a *problemApp) leak(w http.ResponseWriter, r *http.Request) {
	kb := queryInt(r, "kb", 512, 1, 8*1024)
	buffer := make([]byte, kb<<10)
	for page := 0; page < len(buffer); page += 4096 {
		buffer[page] = 1
	}

	a.heapMu.Lock()
	defer a.heapMu.Unlock()
	if a.retained+len(buffer) > a.maxRetainedBytes {
		http.Error(w, "intentional leak safety cap reached; restart the service", http.StatusInsufficientStorage)
		return
	}
	// BUG: a process-lifetime slice keeps every request buffer reachable.
	a.retainedHeap = append(a.retainedHeap, buffer)
	a.retained += len(buffer)
	writeJSON(w, map[string]any{"retained_bytes": a.retained, "objects": len(a.retainedHeap)})
}

func (a *problemApp) slowIO(w http.ResponseWriter, r *http.Request) {
	delay := time.Duration(queryInt(r, "ms", 500, 1, 5_000)) * time.Millisecond
	reader, writer, err := os.Pipe()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	go delayedWrite(writer, delay)
	data, err := readSlowDependency(reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"result": string(data), "waited": delay.String()})
}

//go:noinline
func delayedWrite(writer *os.File, delay time.Duration) {
	defer writer.Close()
	time.Sleep(delay)
	_, _ = writer.Write([]byte("dependency response"))
}

//go:noinline
func readSlowDependency(reader *os.File) ([]byte, error) {
	// BUG: the request has no context-aware timeout or cancellation and waits
	// on a slow dependency. This consumes wall time, not CPU time.
	return io.ReadAll(reader)
}

func (a *problemApp) goroutineLeak(w http.ResponseWriter, r *http.Request) {
	n := queryInt(r, "n", 25, 1, 1_000)
	if a.leakedRoutines.Load()+int64(n) > 10_000 {
		http.Error(w, "intentional goroutine leak safety cap reached; restart the service", http.StatusInsufficientStorage)
		return
	}
	for i := 0; i < n; i++ {
		a.leakedRoutines.Add(1)
		go leakedWorker(a.neverReleased)
	}
	writeJSON(w, map[string]any{"leaked_goroutines": a.leakedRoutines.Load()})
}

//go:noinline
func leakedWorker(neverReleased <-chan struct{}) {
	// BUG: nobody closes or sends on this process-lifetime channel.
	<-neverReleased
}

func (a *problemApp) channelContention(w http.ResponseWriter, r *http.Request) {
	messages := queryInt(r, "messages", 10, 1, 1_000)
	delay := time.Duration(queryInt(r, "consumer_ms", 25, 0, 1_000)) * time.Millisecond
	processed := sendWithBackpressure(messages, delay)
	writeJSON(w, map[string]any{"processed": processed})
}

//go:noinline
func sendWithBackpressure(messages int, consumerDelay time.Duration) int {
	queue := make(chan int) // BUG: no burst capacity and a deliberately slow consumer.
	done := make(chan int, 1)
	go func() {
		count := 0
		for range queue {
			count++
			time.Sleep(consumerDelay)
		}
		done <- count
	}()
	for i := 0; i < messages; i++ {
		queue <- i // block profile attributes delay to this send path.
	}
	close(queue)
	return <-done
}

func (a *problemApp) mutexContention(w http.ResponseWriter, r *http.Request) {
	workers := queryInt(r, "workers", 16, 1, 128)
	loops := queryInt(r, "loops", 20, 1, 1_000)
	hold := time.Duration(queryInt(r, "hold_ms", 2, 0, 100)) * time.Millisecond

	var group sync.WaitGroup
	group.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer group.Done()
			for j := 0; j < loops; j++ {
				a.holdHotLock(hold)
			}
		}()
	}
	group.Wait()
	writeJSON(w, map[string]any{"counter": a.sharedCounter})
}

//go:noinline
func (a *problemApp) holdHotLock(hold time.Duration) {
	a.hotMu.Lock()
	defer a.hotMu.Unlock()
	// BUG: sleeping or doing slow work while holding a global mutex serializes
	// all callers. Mutex profiles point at the holder/unlock stack.
	time.Sleep(hold)
	a.sharedCounter++
}

func (a *problemApp) stats(w http.ResponseWriter, _ *http.Request) {
	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)
	a.heapMu.Lock()
	retained := a.retained
	objects := len(a.retainedHeap)
	a.heapMu.Unlock()
	writeJSON(w, map[string]any{
		"goroutines":               runtime.NumGoroutine(),
		"intentionally_leaked":     a.leakedRoutines.Load(),
		"retained_bytes":           retained,
		"retained_objects":         objects,
		"go_heap_alloc_bytes":      memory.HeapAlloc,
		"go_total_allocated_bytes": memory.TotalAlloc,
	})
}

func queryInt(r *http.Request, name string, fallback, low, high int) int {
	value, err := strconv.Atoi(r.URL.Query().Get(name))
	if err != nil || value < low {
		return fallback
	}
	if value > high {
		return high
	}
	return value
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}

func textHandler(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = io.WriteString(w, body)
	}
}
