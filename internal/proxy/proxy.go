package proxy

import (
	"caching-proxy/internal/cache"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Proxy struct {
	cache  *cache.Cache
	origin *url.URL
	client *http.Client
}

// NewProxy initializes and returns a new Proxy instance with the provided origin and cache.
func NewProxy(origin string, cache *cache.Cache) (*Proxy, error) {
	originURL, err := url.Parse(origin)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		cache:  cache,
		origin: originURL,
		client: &http.Client{},
	}, nil
}

// ServeHTTP handles HTTP requests, checks the cache, and forwards request to the origin if necessary.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := p.generateCacheKey(r)

	// Check and serve from cache if the entry exists
	if entry, exists := p.cache.Get(key); exists {
		p.serveFromCache(w, entry)
		return
	}

	// Forward the request to the origin server
	if err := p.forwardRequest(w, r, key); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// generateCacheKey generates a cache key for the request.
func (p *Proxy) generateCacheKey(r *http.Request) string {
	// Ensure that the cache key is safe by escaping the URL path.
	return r.Method + url.PathEscape(r.URL.Path)
}

// serveFromCache writes the cached response to the ResponseWriter.
func (p *Proxy) serveFromCache(w http.ResponseWriter, entry cache.Entry) {
	for k, v := range entry.Headers {
		w.Header()[k] = v
	}
	w.Header().Set("X-Cache", "HIT")
	w.Header().Set("Content-Type", entry.ContentType)
	w.Write(entry.Data)
}

// forwardRequest forwards the request to the origin server, caches the response, and writes it to the ResponseWriter.
func (p *Proxy) forwardRequest(w http.ResponseWriter, r *http.Request, key string) error {
	targetURL := p.origin.ResolveReference(r.URL).String()
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		return err
	}

	// Copy original headers
	copyHeaders(req.Header, r.Header)

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Cache the response if it is cacheable
	if isResponseCacheable(resp) {
		p.cache.Set(key, body, resp.Header, resp.Header.Get("Content-Type"))
	}

	// Set response headers
	copyHeaders(w.Header(), resp.Header)
	w.Header().Set("X-Cache", "MISS")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	return nil
}

// copyHeaders copies headers from the source to the destination.
func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

// isResponseCacheable determines whether the response is cacheable.
func isResponseCacheable(resp *http.Response) bool {
	cacheControl := strings.ToLower(resp.Header.Get("Cache-Control"))
	return !strings.Contains(cacheControl, "no-store") && !strings.Contains(cacheControl, "private")
}
