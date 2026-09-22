package telemetry

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
	inFlight prometheus.Gauge
	registry *prometheus.Registry
}

func NewMetrics(enabled bool) (*Metrics, error) {
	if !enabled {
		return nil, nil
	}

	requests := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "developers_cv",
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of HTTP requests.",
	}, []string{"method", "route", "status"})
	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "developers_cv",
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "HTTP request duration in seconds.",
	}, []string{"method", "route"})
	inFlight := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "developers_cv",
		Subsystem: "http",
		Name:      "requests_in_flight",
		Help:      "Current number of in-flight HTTP requests.",
	})

	registry := prometheus.NewRegistry()
	if err := registry.Register(requests); err != nil {
		return nil, fmt.Errorf("register HTTP request counter: %w", err)
	}
	if err := registry.Register(duration); err != nil {
		return nil, fmt.Errorf("register HTTP request duration: %w", err)
	}
	if err := registry.Register(inFlight); err != nil {
		return nil, fmt.Errorf("register HTTP in-flight gauge: %w", err)
	}
	return &Metrics{
		requests: requests,
		duration: duration,
		inFlight: inFlight,
		registry: registry,
	}, nil
}

func (m *Metrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		m.inFlight.Inc()
		defer m.inFlight.Dec()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		m.requests.WithLabelValues(c.Request.Method, route, strconv.Itoa(c.Writer.Status())).Inc()
		m.duration.WithLabelValues(c.Request.Method, route).Observe(time.Since(start).Seconds())
	}
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}
