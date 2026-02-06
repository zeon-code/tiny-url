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

		assert.Equal(t, float64(1), fake.MetricBackend.LastTimingRate)
		assert.NotNil(t, fake.MetricBackend.LastTimingDuration)
		assert.Equal(t, "http.request", fake.MetricBackend.LastTimingName)
		assert.Equal(t, []string{"method:GET", "route:/api/v1/url", "status:200"}, fake.MetricBackend.LastTimingTags)
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

		metrics.CacheHit("key", time.Since(time.Now()))

		assert.Equal(t, float64(1), fake.MetricBackend.LastTimingRate)
		assert.NotNil(t, fake.MetricBackend.LastTimingDuration)
		assert.Equal(t, "cache.hit", fake.MetricBackend.LastTimingName)
		assert.Equal(t, []string{"key:key"}, fake.MetricBackend.LastTimingTags)
	})

	t.Run("should send cache hit metric with err", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		fake.MetricBackend.Err = context.Canceled
		metrics.CacheHit("key", time.Since(time.Now()))
	})

	t.Run("should send cache miss metric", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		metrics.CacheMiss("key", time.Since(time.Now()))

		assert.Equal(t, float64(1), fake.MetricBackend.LastTimingRate)
		assert.NotNil(t, fake.MetricBackend.LastTimingDuration)
		assert.Equal(t, "cache.miss", fake.MetricBackend.LastTimingName)
		assert.Equal(t, []string{"key:key"}, fake.MetricBackend.LastTimingTags)
	})

	t.Run("should send cache miss metric with err", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		fake.MetricBackend.Err = context.Canceled
		metrics.CacheMiss("key", time.Since(time.Now()))
	})

	t.Run("should send cache error metric", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		err := context.Canceled
		metrics.CacheError("key", err.Error())

		assert.Equal(t, float64(1), fake.MetricBackend.LastIncrRate)
		assert.Equal(t, "cache.error", fake.MetricBackend.LastIncrName)
		assert.Equal(t, []string{"key:key", "error:" + err.Error()}, fake.MetricBackend.LastIncrTags)
	})

	t.Run("should send cache error metric with err", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		err := context.Canceled
		fake.MetricBackend.Err = err
		metrics.CacheError("key", err.Error())
	})

	t.Run("should send memory invalid metric", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		metrics.MemoryInvalid("key")

		assert.Equal(t, float64(1), fake.MetricBackend.LastIncrRate)
		assert.Equal(t, "memory.invalid", fake.MetricBackend.LastIncrName)
		assert.Equal(t, []string{"key:key"}, fake.MetricBackend.LastIncrTags)
	})

	t.Run("should send memory invalid metric with err", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		fake.MetricBackend.Err = context.Canceled
		metrics.MemoryInvalid("key")
	})

	t.Run("should send memory hit metric", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		metrics.MemoryHit("key", time.Since(time.Now()))

		assert.Equal(t, float64(1), fake.MetricBackend.LastTimingRate)
		assert.NotNil(t, fake.MetricBackend.LastTimingDuration)
		assert.Equal(t, "memory.hit", fake.MetricBackend.LastTimingName)
		assert.Equal(t, []string{"key:key"}, fake.MetricBackend.LastTimingTags)
	})

	t.Run("should send memory hit metric with err", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		fake.MetricBackend.Err = context.Canceled
		metrics.MemoryHit("key", time.Since(time.Now()))
	})

	t.Run("should send memory miss metric", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		metrics.MemoryMiss("key", time.Since(time.Now()))

		assert.Equal(t, float64(1), fake.MetricBackend.LastTimingRate)
		assert.NotNil(t, fake.MetricBackend.LastTimingDuration)
		assert.Equal(t, "memory.miss", fake.MetricBackend.LastTimingName)
		assert.Equal(t, []string{"key:key"}, fake.MetricBackend.LastTimingTags)
	})

	t.Run("should send memory hit metric with err", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		fake.MetricBackend.Err = context.Canceled
		metrics.MemoryMiss("key", time.Since(time.Now()))
	})

	t.Run("should send memory bypass metric", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		metrics.MemoryBypassed()

		assert.Equal(t, float64(1), fake.MetricBackend.LastIncrRate)
		assert.Equal(t, "memory.bypassed", fake.MetricBackend.LastIncrName)
		assert.Equal(t, []string{}, fake.MetricBackend.LastIncrTags)
	})

	t.Run("should send cache bypass metric with err", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		fake.MetricBackend.Err = context.Canceled
		metrics.MemoryBypassed()
	})

	t.Run("should send db query metric", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())
		query := "SELECT * FROM url"

		metrics.DBQuery(query, time.Since(time.Now()))

		assert.Equal(t, float64(1), fake.MetricBackend.LastTimingRate)
		assert.NotNil(t, fake.MetricBackend.LastTimingDuration)
		assert.Equal(t, "db.query", fake.MetricBackend.LastTimingName)
		assert.Equal(t, []string{"query:" + query}, fake.MetricBackend.LastTimingTags)
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
		assert.Equal(t, "db.error", fake.MetricBackend.LastIncrName)
		assert.Equal(t, []string{"query:" + query, "error:" + err.Error()}, fake.MetricBackend.LastIncrTags)
	})

	t.Run("should send db query error metric with err", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		metrics := metric.NewDatadogClient(fake.MetricBackend, "test", fake.Logger())

		err := context.Canceled
		fake.MetricBackend.Err = err
		metrics.DBError("SELECT * FROM url", err.Error())
	})
}
