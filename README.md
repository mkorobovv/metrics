# Metrics example in Go

This project demonstrates how to implement and use Prometheus metrics in a Go HTTP server, including request counting and latency measurement.

## Features

- HTTP request metrics (count and latency)

- Prometheus metrics endpoint

- Example worker that generates traffic

- Graceful shutdown handling

- Structured logging

## Metrics Collected

The following metrics are exposed at `/metrics`:

1. `http_requests_total`
   - **Type**: Counter
   - **Labels**: method, path, code
   - **Description**: Total number of HTTP requests

2. `http_request_latency_seconds`
   - **Type**: Histogram
   - **Labels**: method, path
   - **Description**: Response latency of HTTP requests
   - 
3. Standard Go and process metrics (from Prometheus collectors)

## Getting Started

### Running the Example

1. **Clone the repository:**
```shell
git clone <repository-url>
cd metrics-example
```

2. **Run the server:**

```shell
go run ./examples/main.go
```

3. **In another terminal, generate traffic (optional):**

```shell
while true; do curl http://localhost:8080/health; sleep 0.1; done
```

## Customization
To use these metrics in your own project:

1. **Import:**
```go
import "github.com/mkorobovv/metrics"
```

2. **Create and register metrics:**
```go
serverMetrics := metrics.NewServerMetrics()
registry.MustRegister(serverMetrics)
```

3. **Wrap your handlers:**

```go
r.HandleFunc("/yourpath", serverMetrics.WrapHandlerFunc(yourHandler))
```
