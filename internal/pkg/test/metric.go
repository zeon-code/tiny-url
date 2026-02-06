package test

import (
	"time"
)

type FakeMetric struct {
	// Cache
	LastCacheHit          bool
	LastCacheMiss         bool
	LastMemoryInvalid     bool
	LastCacheErr          string
	LastCacheKeyErr       string
	LastMemoryHitKey      string
	LastMemoryHitLatency  time.Duration
	LastMemoryMissKey     string
	LastMemoryMissLatency time.Duration
	LastMemoryBypass      bool

	// Database
	LastQuery      string
	LastDuration   time.Duration
	LastDBErr      string
	LastDBQueryErr string

	//HTTP
	LastHTTPRequestPath       string
	LastHTTPRequestMethod     string
	LastHTTPRequestStatusCode int
	LastHTTPRequestDuration   time.Duration
}

func (m *FakeMetric) HTTPRequest(method string, path string, statusCode int, duration time.Duration) {
	m.LastHTTPRequestMethod = method
	m.LastHTTPRequestPath = path
	m.LastHTTPRequestStatusCode = statusCode
	m.LastHTTPRequestDuration = duration
}

func (m *FakeMetric) CacheHit(string, time.Duration) {
	m.LastCacheHit = true
}

func (m *FakeMetric) CacheMiss(string, time.Duration) {
	m.LastCacheMiss = true
}

func (m *FakeMetric) CacheError(key string, err string) {
	m.LastCacheKeyErr = key
	m.LastCacheErr = err
}

func (m *FakeMetric) MemoryHit(key string, duration time.Duration) {
	m.LastMemoryHitKey = key
	m.LastMemoryHitLatency = duration
}

func (m *FakeMetric) MemoryMiss(key string, duration time.Duration) {
	m.LastMemoryMissKey = key
	m.LastMemoryMissLatency = duration
}

func (m *FakeMetric) MemoryInvalid(string) {
	m.LastMemoryInvalid = true
}

func (m *FakeMetric) MemoryBypassed() {
	m.LastMemoryBypass = true
}

func (m *FakeMetric) DBQuery(q string, d time.Duration) {
	m.LastQuery = q
	m.LastDuration = d
}

func (m *FakeMetric) DBError(q string, err string) {
	m.LastDBErr = err
	m.LastDBQueryErr = q
}

func NewFakeMetric() *FakeMetric {
	return &FakeMetric{}
}
