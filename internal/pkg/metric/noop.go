package metric

import (
	"log/slog"
	"time"

	"github.com/zeon-code/tiny-url/internal/pkg/config"
)

type NoopClient struct{}

func (NoopClient) HTTPRequest(string, string, int, time.Duration) {}
func (NoopClient) CacheHit(string, time.Duration)                 {}
func (NoopClient) CacheMiss(string, time.Duration)                {}
func (NoopClient) CacheError(string, string)                      {}
func (NoopClient) MemoryHit(string, time.Duration)                {}
func (NoopClient) MemoryMiss(string, time.Duration)               {}
func (NoopClient) MemoryInvalid(string)                           {}
func (NoopClient) MemoryBypassed()                                {}
func (NoopClient) DBQuery(string, time.Duration)                  {}
func (NoopClient) DBError(string, string)                         {}

func NewNoopClient(conf config.MetricConfiguration, logger *slog.Logger) NoopClient {
	return NoopClient{}
}
