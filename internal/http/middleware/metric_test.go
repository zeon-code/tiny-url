package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeon-code/tiny-url/internal/http/middleware"
	"github.com/zeon-code/tiny-url/internal/pkg/test"
)

func TestHTTPMetrics(t *testing.T) {
	t.Run("should record http metrics", func(t *testing.T) {
		metrics := test.NewFakeMetric()

		handler := middleware.HTTPMetrics(metrics)(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
		)

		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

		assert.Equal(t, "/", metrics.LastHTTPRequestPath)
		assert.NotNil(t, metrics.LastHTTPRequestDuration)
		assert.Equal(t, http.MethodGet, metrics.LastHTTPRequestMethod)
		assert.Equal(t, http.StatusOK, metrics.LastHTTPRequestStatusCode)
	})
}
