package proxy

import (
	"caching-proxy/internal/cache"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxy(t *testing.T) {
	// Create a mock origin server
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"test": "data"}`))
	}))
	defer origin.Close()

	// Create cache and proxy
	c := cache.NewCache()
	p, err := NewProxy(origin.URL, c)
	if err != nil {
		t.Fatalf("Failed to create proxy: %v", err)
	}

	// Create test server using our proxy
	proxy := httptest.NewServer(http.HandlerFunc(p.ServeHTTP))
	defer proxy.Close()

	// Test cases
	tests := []struct {
		name           string
		expectedCache  string
		expectedStatus int
	}{
		{"First request should be a cache MISS", "MISS", http.StatusOK},
		{"Second request should be a cache HIT", "HIT", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(proxy.URL + "/test")
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			cacheHeader := resp.Header.Get("X-Cache")
			if cacheHeader != tt.expectedCache {
				t.Errorf("Expected X-Cache: %s, got %s", tt.expectedCache, cacheHeader)
			}
		})
	}
}
