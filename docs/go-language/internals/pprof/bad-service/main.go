// bad-service contains deliberate, bounded performance bugs for pprof study.
// Do not copy these handlers into a real service.
package main

import (
	"flag"
	"log"
	"net/http"
	httppprof "net/http/pprof"
	"runtime"
	"time"
)

func main() {
	appAddr := flag.String("http", "127.0.0.1:8081", "application listen address")
	pprofAddr := flag.String("pprof", "127.0.0.1:6061", "pprof listen address")
	maxRetainedMiB := flag.Int("max-retained-mib", 128, "safety cap for the intentional heap leak")
	blockRate := flag.Int("block-rate", 1, "nanoseconds of blocking per block-profile sample; 0 disables")
	mutexFraction := flag.Int("mutex-fraction", 1, "sample 1/N mutex contentions; 0 disables")
	flag.Parse()

	runtime.SetBlockProfileRate(*blockRate)
	runtime.SetMutexProfileFraction(*mutexFraction)

	go func() {
		log.Printf("pprof UI/index: http://%s/debug/pprof/", *pprofAddr)
		if err := http.ListenAndServe(*pprofAddr, newPprofMux()); err != nil {
			log.Fatalf("pprof server: %v", err)
		}
	}()

	app := newProblemApp(*maxRetainedMiB << 20)
	server := &http.Server{
		Addr:              *appAddr,
		Handler:           app.routes(),
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	log.Printf("bad service: http://%s (open / for scenario links)", *appAddr)
	log.Fatal(server.ListenAndServe())
}

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
