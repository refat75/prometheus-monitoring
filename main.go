package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Prometheus Metrics
var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests.",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "HTTP response status code.",
	},
	[]string{"status"},
)

var httpDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "HTTP request latencies in seconds.",
	},
	[]string{"path"},
)

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routePattern := chi.RouteContext(r.Context()).RoutePattern()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(routePattern))
		rw := NewResponseWriter(w)

		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode
		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequests.WithLabelValues(routePattern).Inc()

		timer.ObserveDuration()
	})
}

func init() {
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(responseStatus)
	prometheus.MustRegister(httpDuration)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Register Middleware
	r.Use(prometheusMiddleware)
	r.Handle("/prometheus", promhttp.Handler())

	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "./static/index.html")
	})

	fmt.Println("Serving requests on port 8080")
	http.ListenAndServe(":8080", r)
}
