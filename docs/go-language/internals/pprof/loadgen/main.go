// loadgen repeatedly calls one URL so a performance problem has enough samples
// to become visible. It intentionally uses only the standard library.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	target := flag.String("url", "http://127.0.0.1:8080/work?size=2000", "full URL to request")
	duration := flag.Duration("duration", 15*time.Second, "load duration")
	concurrency := flag.Int("concurrency", 8, "number of request workers")
	timeout := flag.Duration("timeout", 10*time.Second, "per-request timeout")
	flag.Parse()

	if *concurrency < 1 || *duration <= 0 {
		fmt.Fprintln(os.Stderr, "concurrency and duration must be positive")
		os.Exit(2)
	}

	transport := &http.Transport{
		MaxIdleConns:        *concurrency * 2,
		MaxIdleConnsPerHost: *concurrency,
		IdleConnTimeout:     30 * time.Second,
	}
	client := &http.Client{Transport: transport, Timeout: *timeout}
	defer transport.CloseIdleConnections()

	ctx, cancel := context.WithTimeout(context.Background(), *duration)
	defer cancel()
	start := time.Now()
	var requests atomic.Int64
	var failures atomic.Int64
	var totalLatency atomic.Int64
	var maxLatency atomic.Int64

	var workers sync.WaitGroup
	workers.Add(*concurrency)
	for i := 0; i < *concurrency; i++ {
		go func() {
			defer workers.Done()
			for ctx.Err() == nil {
				requestStart := time.Now()
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, *target, nil)
				if err != nil {
					failures.Add(1)
					return
				}
				response, err := client.Do(req)
				latency := time.Since(requestStart).Nanoseconds()
				totalLatency.Add(latency)
				updateMax(&maxLatency, latency)
				requests.Add(1)
				if err != nil {
					if ctx.Err() == nil {
						failures.Add(1)
					}
					continue
				}
				_, copyErr := io.Copy(io.Discard, response.Body)
				closeErr := response.Body.Close()
				if response.StatusCode >= 400 || copyErr != nil || closeErr != nil {
					failures.Add(1)
				}
			}
		}()
	}
	workers.Wait()

	elapsed := time.Since(start)
	count := requests.Load()
	average := time.Duration(0)
	if count > 0 {
		average = time.Duration(totalLatency.Load() / count)
	}
	fmt.Printf("url=%s\n", *target)
	fmt.Printf("elapsed=%s requests=%d failures=%d rate=%.1f req/s avg=%s max=%s\n",
		elapsed.Round(time.Millisecond), count, failures.Load(), float64(count)/elapsed.Seconds(),
		average.Round(time.Microsecond), time.Duration(maxLatency.Load()).Round(time.Microsecond))
}

func updateMax(maximum *atomic.Int64, value int64) {
	for {
		current := maximum.Load()
		if value <= current || maximum.CompareAndSwap(current, value) {
			return
		}
	}
}
