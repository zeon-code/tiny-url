package middleware

import (
	"net/http"
	"time"

	"github.com/zeon-code/tiny-url/internal/pkg/metric"
)

type HTTPMetricResponse struct {
	StatusCode int
	Response   http.ResponseWriter
}

func (r *HTTPMetricResponse) Header() http.Header {
	return r.Response.Header()
}

func (r *HTTPMetricResponse) Write(data []byte) (int, error) {
	return r.Response.Write(data)
}

func (r *HTTPMetricResponse) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.Response.WriteHeader(statusCode)
}

func HTTPMetrics(metrics metric.MetricClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()

			resp := &HTTPMetricResponse{Response: w}
			next.ServeHTTP(resp, r)

			metrics.HTTPRequest(r.Method, r.URL.Path, resp.StatusCode, time.Since(now))
		})
	}
}
