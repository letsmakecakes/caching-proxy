package cache

import "testing"

// setupCache initializes a cache and returns it along with test data.
func setupCache() (*Cache, string, []byte, map[string][]string, string) {
	c := NewCache()
	key := "test-key"
	data := []byte("test-data")
	headers := map[string][]string{
		"Content-Type": {"application/json"},
	}
	contentType := "application/json"
	return c, key, data, headers, contentType
}

func TestCache(t *testing.T) {
	c, key, data, headers, contentType := setupCache()

	t.Run("Set and Get", func(t *testing.T) {
		c.Set(key, data, headers, contentType)
		entry, exists := c.Get(key)

		if !exists {
			t.Error("Expected cache entry to exist")
		}

		if string(entry.Data) != string(data) {
			t.Errorf("Expected content type %s, got %s", contentType, entry.ContentType)
		}

		for k, v := range headers {
			if entry.Headers[k][0] != v[0] {
				t.Errorf("Expected header %s to be %s, got %s", k, v[0], entry.Headers[k][0])
			}
		}
	})
}
