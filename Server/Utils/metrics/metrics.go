package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"service", "route", "method", "status", "version"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets, // 生产可以按你们SLO自定义桶
		},
		[]string{"service", "route", "method", "status", "version"},
	)

	InFlight = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_in_flight_requests",
			Help: "Current number of in-flight HTTP requests",
		},
		[]string{"service", "route", "version"},
	)

	// ====== 可选补充：下游调用（peer 维度）======
	OutboundDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "outbound_request_duration_seconds",
			Help:    "Outbound request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "peer_service", "result", "version"},
	)

	OutboundTimeouts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "outbound_timeouts_total",
			Help: "Total outbound timeouts",
		},
		[]string{"service", "peer_service", "version"},
	)

	RetriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "outbound_retries_total",
			Help: "Total outbound retries",
		},
		[]string{"service", "peer_service", "version"},
	)

	FallbackTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "fallback_total",
			Help: "Total number of fallbacks",
		},
		[]string{"service", "route", "version"},
	)

	RateLimitedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limited_total",
			Help: "Total number of rate-limited requests",
		},
		[]string{"service", "route", "version"},
	)
)

func MustRegister() {
	prometheus.MustRegister(
		RequestsTotal,
		RequestDuration,
		InFlight,
		OutboundDuration,
		OutboundTimeouts,
		RetriesTotal,
		FallbackTotal,
		RateLimitedTotal,
	)
}

func ObserveHTTP(service, route, method string, status int, version string, start time.Time) {
	st := strconv.Itoa(status)
	RequestsTotal.WithLabelValues(service, route, method, st, version).Inc()
	RequestDuration.WithLabelValues(service, route, method, st, version).Observe(time.Since(start).Seconds())
}
