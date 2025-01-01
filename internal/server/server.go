package server

import (
	"caching-proxy/internal/cache"
	"caching-proxy/internal/config"
	"caching-proxy/internal/proxy"
	"fmt"
	"net/http"
)

// Server represents the caching proxy server.
type Server struct {
	cfg   *config.Config
	cache *cache.Cache
}

// NewServer initializes and returns a new Server instance.
func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:   cfg,
		cache: cache.NewCache(),
	}
}

// Start begins serving HTTP request using the proxy handler.
func (s *Server) Start() error {
	proxyHandler, err := proxy.NewProxy(s.cfg.Origin, s.cache)
	if err != nil {
		return fmt.Errorf("failed to create proxy: %w", err)
	}

	address := fmt.Sprintf(":%d", s.cfg.Port)
	fmt.Printf("Starting server on %s\n", address)
	return http.ListenAndServe(address, proxyHandler)
}

// ClearCache clears the global cache instance.
func (s *Server) ClearCache() {
	s.cache.Clear()
}
