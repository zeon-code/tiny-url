package metric

import "time"

type NoopMetrics struct{}

func (NoopMetrics) HTTPRequest(string, string, int, time.Duration) {}
func (NoopMetrics) CacheHit(string)                                {}
func (NoopMetrics) CacheMiss(string)                               {}
func (NoopMetrics) CacheInvalid(string)                            {}
func (NoopMetrics) CacheError(string, string)                      {}
func (NoopMetrics) CacheLatency(string, time.Duration)             {}
func (NoopMetrics) CacheBypassed()                                 {}
func (NoopMetrics) DBQuery(string, time.Duration)                  {}
func (NoopMetrics) DBError(string, string)                         {}
