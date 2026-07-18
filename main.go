package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	addr := flag.String("addr", envOr("OMEGA_ADDR", ":8080"), "address to listen on")
	docsDir := flag.String("docs", envOr("OMEGA_DOCS", "docs"), "directory containing Markdown documents")
	flag.Parse()

	site, err := newSite(os.DirFS(*docsDir))
	if err != nil {
		log.Fatalf("load docs: %v", err)
	}

	server := &http.Server{
		Addr:              *addr,
		Handler:           site.routes(),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	displayAddr := *addr
	if strings.HasPrefix(displayAddr, ":") {
		displayAddr = "localhost" + displayAddr
	}
	log.Printf("Omega Learn is serving %d articles at http://%s", len(site.articles), displayAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
