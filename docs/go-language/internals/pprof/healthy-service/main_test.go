package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	recorder := httptest.NewRecorder()
	newAppMux().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if recorder.Code != http.StatusOK || recorder.Body.String() != "ok\n" {
		t.Fatalf("unexpected response: code=%d body=%q", recorder.Code, recorder.Body.String())
	}
}

func TestWorkIsBoundedAndReturnsJSON(t *testing.T) {
	recorder := httptest.NewRecorder()
	newAppMux().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/work?size=999999", nil))

	var result workResult
	if err := json.NewDecoder(recorder.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.Items != 50_000 {
		t.Fatalf("items=%d, want upper bound 50000", result.Items)
	}
}

func TestPprofIsNotOnApplicationMux(t *testing.T) {
	recorder := httptest.NewRecorder()
	newAppMux().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil))
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("pprof leaked onto app listener: code=%d", recorder.Code)
	}
}
