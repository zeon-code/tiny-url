package metric

import (
	"log/slog"
	"time"

	"github.com/zeon-code/tiny-url/internal/pkg/config"
)

type NoopClient struct{}

func (NoopClient) HTTPRequest(string, string, int, time.Duration) {}
func (NoopClient) CacheHit(string)                                {}
func (NoopClient) CacheMiss(string)                               {}
func (NoopClient) CacheInvalid(string)                            {}
func (NoopClient) CacheError(string, string)                      {}
func (NoopClient) CacheLatency(string, time.Duration)             {}
func (NoopClient) CacheBypassed()                                 {}
func (NoopClient) DBQuery(string, time.Duration)                  {}
func (NoopClient) DBError(string, string)                         {}

func NewNoopClient(conf config.MetricConfiguration, logger *slog.Logger) NoopClient {
	return NoopClient{}
}
