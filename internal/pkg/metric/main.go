package metric

import (
	"log/slog"
	"strings"
	"time"

	"github.com/zeon-code/tiny-url/internal/pkg/config"
)

// Metric defines a vendor-agnostic interface for emitting
// application-level observability signals.
type MetricClient interface {

	// HTTPRequest records the completion of an HTTP request.
	// It should be called once per request, after the response
	// status is known, and must include request duration.
	HTTPRequest(string, string, int, time.Duration)

	// CacheHit records a successful cache lookup where the
	// requested value was found and used.
	CacheHit(string)

	// CacheMiss records a cache lookup where no value was found
	// and a fallback (e.g. database) was required.
	CacheMiss(string)

	// CacheInvalid records a cache entry that existed but could
	// not be used (e.g. stale, malformed, or failed validation).
	CacheInvalid(string)

	// CacheErr records a cache operation that failed due to
	// backend or connectivity errors.
	CacheError(string, string)

	// CacheLatency records the duration taken by cache operations such as reads,
	// writes, deletions, or full read-through paths.
	CacheLatency(string, time.Duration)

	// CacheBypassed records that cache logic was intentionally skipped.
	CacheBypassed()

	// DBQuery records the execution of a database query.
	// It should include the logical query name and execution time.
	DBQuery(string, time.Duration)

	// DBError records a database operation that failed.
	// This should be used for query errors, timeouts, or
	// connection failures.
	DBError(string, string)
}

func NewMetricClient(conf config.Configuration, logger *slog.Logger) MetricClient {
	integration, _ := conf.Metric().Integration()

	switch strings.ToLower(integration) {
	case "datadog":
		client, err := NewDatadogClientFromConf(conf.Metric(), logger)

		if err != nil {
			logger.Error("error building datadog client, returning fallback NoopClient", slog.Any("error", err))
			return NewNoopClient(conf.Metric(), logger)
		}

		return client
	default:
		return NewNoopClient(conf.Metric(), logger)
	}
}
