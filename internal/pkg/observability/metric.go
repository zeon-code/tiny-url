package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/metric"
)

// Metric defines a vendor-agnostic interface for emitting
// application-level observability signals.
type MetricClient interface {
	// MemoryHit records the duration to resolve a value when it is found in memory cache.
	// The duration includes cache lookup and value deserialization, but excludes fallback logic.
	MemoryHit(context.Context, string, time.Duration)

	// MemoryMiss records the duration to resolve a value when it is not found in memory cache.
	// The duration includes cache lookup, fallback data fetch, value serialization,
	// and cache population.
	MemoryMiss(context.Context, string, time.Duration)

	// MemoryInvalid records a cache entry that existed but could
	// not be used (e.g. stale, malformed, or failed validation).
	MemoryInvalid(context.Context, string)

	// CacheBypassed records that cache logic was intentionally skipped.
	MemoryBypassed(context.Context)
}

type OtelMetricClient struct {
	meter metric.Meter

	memoryHitLatency   metric.Float64Histogram
	memoryMissLatency  metric.Float64Histogram
	memoryInvalidCount metric.Int64Counter
	memoryBypassCount  metric.Int64Counter
}

func NewMetricClient(meter metric.Meter) (*OtelMetricClient, error) {
	var err error
	client := &OtelMetricClient{}

	client.memoryInvalidCount, err = meter.Int64Counter("tiny_url.memory.invalid.count")

	if err != nil {
		return nil, err
	}

	client.memoryBypassCount, err = meter.Int64Counter("tiny_url.memory.bypassed.count")

	if err != nil {
		return nil, err
	}

	client.memoryHitLatency, err = meter.Float64Histogram(
		"tiny_url.memory.hit.latency",
		metric.WithUnit("ms"),
		metric.WithDescription("Memory hit operation latency"),
	)

	if err != nil {
		return nil, err
	}

	client.memoryMissLatency, err = meter.Float64Histogram(
		"tiny_url.memory.miss.latency",
		metric.WithUnit("ms"),
		metric.WithDescription("Memory miss operation latency"),
	)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (m *OtelMetricClient) MemoryHit(ctx context.Context, name string, d time.Duration) {
	m.memoryHitLatency.Record(
		ctx,
		float64(d.Milliseconds()),
	)
}

func (m *OtelMetricClient) MemoryMiss(ctx context.Context, name string, d time.Duration) {
	m.memoryMissLatency.Record(
		ctx,
		float64(d.Milliseconds()),
	)
}

func (m *OtelMetricClient) MemoryInvalid(ctx context.Context, name string) {
	m.memoryInvalidCount.Add(
		ctx,
		1,
	)
}

func (m *OtelMetricClient) MemoryBypassed(ctx context.Context) {
	m.memoryBypassCount.Add(
		ctx,
		1,
	)
}
