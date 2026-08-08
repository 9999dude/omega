// healthy-service is a small HTTP service with pprof on a separate listener.
// It is the baseline used before investigating deliberately bad examples.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	httppprof "net/http/pprof"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	appAddr := flag.String("http", "127.0.0.1:8080", "application listen address")
	pprofAddr := flag.String("pprof", "127.0.0.1:6060", "pprof listen address")
	blockRate := flag.Int("block-rate", 1, "nanoseconds of blocking per block-profile sample; 0 disables")
	mutexFraction := flag.Int("mutex-fraction", 1, "sample 1/N mutex contentions; 0 disables")
	flag.Parse()

	// CPU, heap, allocs, goroutine, and threadcreate profiles need no switch.
	// Block and mutex profiles are opt-in because collecting them costs work.
	runtime.SetBlockProfileRate(*blockRate)
	runtime.SetMutexProfileFraction(*mutexFraction)

	go func() {
		log.Printf("pprof UI/index: http://%s/debug/pprof/", *pprofAddr)
		if err := http.ListenAndServe(*pprofAddr, newPprofMux()); err != nil {
			log.Fatalf("pprof server: %v", err)
		}
	}()

	server := &http.Server{
		Addr:              *appAddr,
		Handler:           newAppMux(),
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	log.Printf("healthy service: http://%s (try /healthz and /work?size=2000)", *appAddr)
	log.Fatal(server.ListenAndServe())
}

// newPprofMux registers pprof explicitly instead of importing net/http/pprof
// for side effects on http.DefaultServeMux. This makes it harder to
// accidentally expose profiling data on the public application listener.
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

func newAppMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("GET /work", workHandler)
	return mux
}

type workResult struct {
	Items    int    `json:"items"`
	Checksum uint64 `json:"checksum"`
	Preview  string `json:"preview"`
}

func workHandler(w http.ResponseWriter, r *http.Request) {
	size := boundedInt(r, "size", 2_000, 1, 50_000)
	result := buildReport(size)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

// buildReport does enough bounded work to appear in a profile without being a
// performance bug. It preallocates the builder rather than repeatedly growing
// strings, and it does not retain the result after the request ends.
func buildReport(size int) workResult {
	var builder strings.Builder
	builder.Grow(min(size*5, 64*1024))
	var checksum uint64
	for i := 0; i < size; i++ {
		value := uint64((i*31 + 17) % 10_007)
		checksum = checksum*33 ^ value
		if i < 100 {
			_, _ = fmt.Fprintf(&builder, "%d,", value)
		}
	}
	return workResult{Items: size, Checksum: checksum, Preview: builder.String()}
}

func boundedInt(r *http.Request, name string, fallback, low, high int) int {
	value, err := strconv.Atoi(r.URL.Query().Get(name))
	if err != nil || value < low {
		return fallback
	}
	return min(value, high)
}
