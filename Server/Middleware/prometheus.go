package Middleware

import (
	"AITranslatio/Utils/metrics"
	"github.com/gin-gonic/gin"
	"time"
)

func Prometheus(service, version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		route := c.FullPath()
		if route == "" {
			route = "unknown"
		}
		metrics.InFlight.WithLabelValues(service, route, version).Inc()
		defer metrics.InFlight.WithLabelValues(service, route, version).Dec()

		c.Next()

		status := c.Writer.Status()
		metrics.ObserveHTTP(service, route, c.Request.Method, status, version, start)
	}
}
