package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCPUDoesWork(t *testing.T) {
	iterations, _ := burnCPU(time.Millisecond)
	if iterations == 0 {
		t.Fatal("burnCPU completed no iterations")
	}
}

func TestLeakHasSafetyCap(t *testing.T) {
	app := newProblemApp(1024)
	handler := app.routes()

	first := httptest.NewRecorder()
	handler.ServeHTTP(first, httptest.NewRequest(http.MethodGet, "/leak?kb=1", nil))
	if first.Code != http.StatusOK {
		t.Fatalf("first leak response=%d", first.Code)
	}

	second := httptest.NewRecorder()
	handler.ServeHTTP(second, httptest.NewRequest(http.MethodGet, "/leak?kb=1", nil))
	if second.Code != http.StatusInsufficientStorage {
		t.Fatalf("second leak response=%d, want %d", second.Code, http.StatusInsufficientStorage)
	}
}

func TestChannelScenarioCompletes(t *testing.T) {
	if got := sendWithBackpressure(5, 0); got != 5 {
		t.Fatalf("processed=%d, want 5", got)
	}
}

func TestStatsReturnsJSON(t *testing.T) {
	recorder := httptest.NewRecorder()
	newProblemApp(1024).routes().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/stats", nil))
	var result map[string]any
	if err := json.NewDecoder(recorder.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
}
