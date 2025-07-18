package metrics

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type ServerMetrics struct {
	requests *prometheus.CounterVec
	latency  *prometheus.HistogramVec
}

func NewServerMetrics() *ServerMetrics {
	return &ServerMetrics{
		requests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests.",
			},
			[]string{"method", "path", "code"},
		),
		latency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "http_request_latency_seconds",
				Help: "Histogram of response latency of HTTP requests that had been application-level handled by the server.",
			},
			[]string{"method", "path"},
		),
	}
}

func (s *ServerMetrics) Describe(ch chan<- *prometheus.Desc) {
	s.requests.Describe(ch)
	if s.latency != nil {
		s.latency.Describe(ch)
	}
}

func (s *ServerMetrics) Collect(ch chan<- prometheus.Metric) {
	s.requests.Collect(ch)
	if s.latency != nil {
		s.latency.Collect(ch)
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (s *ServerMetrics) WrapHandler(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.URL.Path

		timer := prometheus.NewTimer(s.latency.WithLabelValues(method, path))
		defer timer.ObserveDuration()

		sw := &statusWriter{ResponseWriter: w}

		handler.ServeHTTP(sw, r)

		status := sw.status
		if sw.status == 0 {
			status = http.StatusOK
		}

		s.requests.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
	}
}
