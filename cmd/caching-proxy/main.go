package main

import (
	"caching-proxy/internal/config"
	"caching-proxy/internal/server"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
)

func main() {
	cfg := &config.Config{}

	// Define flags
	flag.IntVar(&cfg.Port, "port", 3000, "Port to run the proxy server on")
	flag.StringVar(&cfg.Origin, "origin", "", "Origin server URL")
	clearCache := flag.Bool("clear-cache", false, "Clear the cache")
	flag.Parse()

	// Validate the Origin URL
	if cfg.Origin == "" {
		log.Fatal("Origin URL is required")
	}
	if _, err := url.ParseRequestURI(cfg.Origin); err != nil {
		log.Fatalf("Invalid Origin URL: %v", err)
	}

	// Create server instance
	srv := server.NewServer(cfg)

	// Clear cache if the flag is set
	if *clearCache {
		srv.ClearCache()
		fmt.Println("Cache is cleared successfully")
		os.Exit(0)
	}

	// Start the server
	log.Printf("Starting proxy server on port %d, forwarding to %s", cfg.Port, cfg.Origin)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
