package metric_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeon-code/tiny-url/internal/pkg/metric"
	"github.com/zeon-code/tiny-url/internal/pkg/test"
)

func TestDatadogClient(t *testing.T) {
	t.Run("should send http metrics", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		metrics.HTTPRequest("GET", "/api/v1/url", 200, time.Since(time.Now()))

		assert.Equal(t, float64(1), fake.MetricBackend.LastIncrRate)
		assert.Equal(t, "tiny_url.http.request.count", fake.MetricBackend.LastIncrName)
		assert.Equal(t, []string{"tiny_url", "env:test", "method:GET", "route:/api/v1/url", "status:200"}, fake.MetricBackend.LastIncrTags)

		assert.Equal(t, float64(1), fake.MetricBackend.LastTimingRate)
		assert.NotNil(t, fake.MetricBackend.LastTimingDuration)
		assert.Equal(t, "tiny_url.http.request.duration", fake.MetricBackend.LastTimingName)
		assert.Equal(t, []string{"tiny_url", "env:test", "method:GET", "route:/api/v1/url", "status:200"}, fake.MetricBackend.LastTimingTags)
	})

	t.Run("should send http metrics with err", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		fake.MetricBackend.Err = context.Canceled
		metrics.HTTPRequest("GET", "/api/v1/url", 200, time.Since(time.Now()))
	})

	t.Run("should send cache hit metric", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		metrics.CacheHit("key")

		assert.Equal(t, float64(1), fake.MetricBackend.LastIncrRate)
		assert.Equal(t, "tiny_url.cache.hit", fake.MetricBackend.LastIncrName)
		assert.Equal(t, []string{"tiny_url", "env:test", "key:key"}, fake.MetricBackend.LastIncrTags)
	})

	t.Run("should send cache hit metric with err", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		fake.MetricBackend.Err = context.Canceled
		metrics.CacheHit("key")
	})

	t.Run("should send cache miss metric", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		metrics.CacheMiss("key")

		assert.Equal(t, float64(1), fake.MetricBackend.LastIncrRate)
		assert.Equal(t, "tiny_url.cache.miss", fake.MetricBackend.LastIncrName)
		assert.Equal(t, []string{"tiny_url", "env:test", "key:key"}, fake.MetricBackend.LastIncrTags)
	})

	t.Run("should send cache miss metric with err", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		fake.MetricBackend.Err = context.Canceled
		metrics.CacheMiss("key")
	})

	t.Run("should send cache invalid metric", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		metrics.CacheInvalid("key")

		assert.Equal(t, float64(1), fake.MetricBackend.LastIncrRate)
		assert.Equal(t, "tiny_url.cache.invalid", fake.MetricBackend.LastIncrName)
		assert.Equal(t, []string{"tiny_url", "env:test", "key:key"}, fake.MetricBackend.LastIncrTags)
	})

	t.Run("should send cache invalid metric with err", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		fake.MetricBackend.Err = context.Canceled
		metrics.CacheInvalid("key")
	})

	t.Run("should send cache error metric", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		err := context.Canceled
		metrics.CacheError("key", err.Error())

		assert.Equal(t, float64(1), fake.MetricBackend.LastIncrRate)
		assert.Equal(t, "tiny_url.cache.error", fake.MetricBackend.LastIncrName)
		assert.Equal(t, []string{"tiny_url", "env:test", "key:key", "error:" + err.Error()}, fake.MetricBackend.LastIncrTags)
	})

	t.Run("should send cache error metric with err", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		err := context.Canceled
		fake.MetricBackend.Err = err
		metrics.CacheError("key", err.Error())
	})

	t.Run("should send cache invalid metric", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		metrics.CacheLatency("key", time.Since(time.Now()))

		assert.Equal(t, float64(1), fake.MetricBackend.LastTimingRate)
		assert.NotNil(t, fake.MetricBackend.LastTimingDuration)
		assert.Equal(t, "tiny_url.cache.latency", fake.MetricBackend.LastTimingName)
		assert.Equal(t, []string{"tiny_url", "env:test", "key:key"}, fake.MetricBackend.LastTimingTags)
	})

	t.Run("should send cache invalid metric with err", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		fake.MetricBackend.Err = context.Canceled
		metrics.CacheLatency("key", time.Since(time.Now()))
	})

	t.Run("should send cache bypass metric", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		metrics.CacheBypassed()

		assert.Equal(t, float64(1), fake.MetricBackend.LastIncrRate)
		assert.Equal(t, "tiny_url.cache.bypassed", fake.MetricBackend.LastIncrName)
		assert.Equal(t, []string{"tiny_url", "env:test"}, fake.MetricBackend.LastIncrTags)
	})

	t.Run("should send cache bypass metric with err", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		fake.MetricBackend.Err = context.Canceled
		metrics.CacheBypassed()
	})

	t.Run("should send db query metric", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())
		query := "SELECT * FROM url"

		metrics.DBQuery(query, time.Since(time.Now()))

		assert.Equal(t, float64(1), fake.MetricBackend.LastTimingRate)
		assert.NotNil(t, fake.MetricBackend.LastTimingDuration)
		assert.Equal(t, "tiny_url.db.query.duration", fake.MetricBackend.LastTimingName)
		assert.Equal(t, []string{"tiny_url", "env:test", "query:" + query}, fake.MetricBackend.LastTimingTags)
	})

	t.Run("should send db query metric with err", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		fake.MetricBackend.Err = context.Canceled
		metrics.DBQuery("key", time.Since(time.Now()))
	})

	t.Run("should send db query error metric", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())
		query := "SELECT * FROM url"
		err := context.Canceled

		metrics.DBError(query, err.Error())

		assert.Equal(t, float64(1), fake.MetricBackend.LastIncrRate)
		assert.Equal(t, "tiny_url.db.error", fake.MetricBackend.LastIncrName)
		assert.Equal(t, []string{"tiny_url", "env:test", "query:" + query, "error:" + err.Error()}, fake.MetricBackend.LastIncrTags)
	})

	t.Run("should send db query error metric with err", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		err := context.Canceled
		fake.MetricBackend.Err = err
		metrics.DBError("SELECT * FROM url", err.Error())
	})
}
