# Caching Proxy

A lightweight caching proxy server written in Go that caches responses from origin servers. It can be used to reduce load on the origin server and improve response times for repeated requests.

## Features

- Forward requests to origin server
- Cache responses for subsequent requests
- Cache hit/miss indicators via headers
- CLI tool with configuration options
- Cache clearing functionality
- Thread-safe implementation

## Installation

### Prerequisites

- Go 1.21 or higher
- Git

### Building from source

1. Clone the repository
```bash
git clone https://github.com/letsmakecakes/caching-proxy.git
cd caching-proxy
```

2. Build the binary
```bash
go build -o caching-proxy cmd/caching-proxy/main.go
```

## Usage

### Starting the proxy server

The basic syntax for starting the proxy server is:

```bash
./caching-proxy --port <port_number> --origin <origin_url>
```

Example:
```bash
./caching-proxy --port 3000 --origin http://dummyjson.com
```

This will start the proxy server on port 3000 and forward requests to http://dummyjson.com.

### Command line flags

- `--port`: Port number for the proxy server (default: 3000)
- `--origin`: Origin server URL (required)
- `--clear-cache`: Clear the cached responses

### Clearing the cache

To clear the cache:
```bash
./caching-proxy --clear-cache
```

## Testing

### Manual Testing

1. Start the proxy server:
```bash
./caching-proxy --port 3000 --origin http://dummyjson.com
```

2. Make a request to the proxy server:
```bash
# First request (cache miss)
curl -v http://localhost:3000/products/1

# Second request (cache hit)
curl -v http://localhost:3000/products/1
```

3. Check the `X-Cache` header in the response:
- `X-Cache: MISS` indicates the response came from the origin server
- `X-Cache: HIT` indicates the response came from the cache

## Performance Testing

You can use Apache Benchmark (ab) to test the performance improvement with caching:

```bash
# Install Apache Benchmark
# On Ubuntu/Debian:
sudo apt-get install apache2-utils
# On macOS:
brew install apache2

# Test without cache (first request)
ab -n 100 -c 10 http://localhost:3000/products/1

# Test with cache (second request)
ab -n 100 -c 10 http://localhost:3000/products/1
```

Compare the response times between cached and non-cached requests.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.