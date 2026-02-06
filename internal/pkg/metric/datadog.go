package metric

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/zeon-code/tiny-url/internal/pkg/config"
)

type DatadogBackend interface {
	Incr(name string, tags []string, rate float64) error
	Timing(name string, value time.Duration, tags []string, rate float64) error
}

type DatadogClient struct {
	client    DatadogBackend
	logger    *slog.Logger
	namespace string
	tags      []string
}

func NewDatadogClientFromConf(conf config.MetricConfiguration, logger *slog.Logger) (*DatadogClient, error) {
	env, err := conf.Environment()

	if err != nil {
		return nil, err
	}

	addr, err := conf.Host()

	if err != nil {
		return nil, err
	}

	port, err := conf.Port()

	if err != nil {
		return nil, err
	}

	client, err := statsd.New(
		fmt.Sprintf("%s:%d", addr, port),
		statsd.WithNamespace("tiny_url."),
		statsd.WithTags([]string{
			"service:tiny_url",
			fmt.Sprintf("env:%s", env),
		}),
	)

	if err != nil {
		return nil, err
	}

	logger.Info("Starting datadog client")
	return NewDatadogClient(client, env, logger), nil
}

func NewDatadogClient(client DatadogBackend, env string, logger *slog.Logger) *DatadogClient {
	return &DatadogClient{
		client:    client,
		logger:    logger,
		namespace: "tiny_url",
		tags:      []string{},
	}
}

func (d *DatadogClient) HTTPRequest(method string, route string, status int, duration time.Duration) {
	tags := append(d.tags, "method:"+method, "route:"+route, "status:"+strconv.Itoa(status))

	if ddErr := d.client.Incr("http.request.count", tags, 1); ddErr != nil {
		d.logger.Error("error while sending request count metric", slog.Any("error", ddErr))
	}

	if ddErr := d.client.Timing("http.request.latency", duration, tags, 1); ddErr != nil {
		d.logger.Error("error while sending request duration metric", slog.Any("error", ddErr))
	}
}

func (d *DatadogClient) CacheHit(key string) {
	if ddErr := d.client.Incr("cache.hit", append(d.tags, "key:"+key), 1); ddErr != nil {
		d.logger.Error("error while sending cache hit metric", slog.Any("error", ddErr))
	}
}

func (d *DatadogClient) CacheMiss(key string) {
	if ddErr := d.client.Incr("cache.miss", append(d.tags, "key:"+key), 1); ddErr != nil {
		d.logger.Error("error while sending cache miss metric", slog.Any("error", ddErr))
	}
}

func (d *DatadogClient) CacheInvalid(key string) {
	if ddErr := d.client.Incr("cache.invalid", append(d.tags, "key:"+key), 1); ddErr != nil {
		d.logger.Error("error while sending cache invalid metric", slog.Any("error", ddErr))
	}
}

func (d *DatadogClient) CacheError(key string, err string) {
	if ddErr := d.client.Incr("cache.error", append(d.tags, "key:"+key, "error:"+err), 1); ddErr != nil {
		d.logger.Error("error while sending cache error metric", slog.Any("error", ddErr))
	}
}

func (d *DatadogClient) CacheLatency(key string, duration time.Duration) {
	if ddErr := d.client.Timing("cache.latency", duration, append(d.tags, "key:"+key), 1); ddErr != nil {
		d.logger.Error("error while sending cache latency metric", slog.Any("error", ddErr))
	}
}

func (d *DatadogClient) CacheBypassed() {
	if ddErr := d.client.Incr("cache.bypassed", d.tags, 1); ddErr != nil {
		d.logger.Error("error while sending cache bypass metric", slog.Any("error", ddErr))
	}
}

func (d *DatadogClient) DBQuery(query string, duration time.Duration) {
	if ddErr := d.client.Timing("db.query.latency", duration, append(d.tags, "query:"+query), 1); ddErr != nil {
		d.logger.Error("error while sending query metric", slog.Any("error", ddErr))
	}
}

func (d *DatadogClient) DBError(query string, err string) {
	if ddErr := d.client.Incr("db.error", append(d.tags, "query:"+query, "error:"+err), 1); ddErr != nil {
		d.logger.Error("error while sending db error metrics", slog.Any("error", ddErr))
	}
}
