package mw

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// MetricsMW — middleware для сбора метрик HTTP-запросов Prometheus.
type MetricsMW struct {
	requestsTotal    *prometheus.CounterVec
	requestDurations *prometheus.HistogramVec
}

func NewMetricsMW() *MetricsMW {
	return &MetricsMW{
		requestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "code"},
		),
		requestDurations: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
	}
}

// Handler — HTTP mw, собирающий метрики по каждому запросу.
func (m *MetricsMW) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method
		start := time.Now()

		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rr, r)

		duration := time.Since(start).Seconds()

		m.requestsTotal.WithLabelValues(method, path, http.StatusText(rr.statusCode)).Inc()
		m.requestDurations.WithLabelValues(method, path).Observe(duration)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (m *MetricsMW) Register(reg prometheus.Registerer) {
	reg.MustRegister(m.requestsTotal)
	reg.MustRegister(m.requestDurations)
}
