package test

import (
	"time"
)

type FakeMetric struct {
	// Cache
	LastCacheHit        bool
	LastCacheMiss       bool
	LastCacheInvalid    bool
	LastCacheErr        string
	LastCacheKeyErr     string
	LastCacheLatency    time.Duration
	LastCacheKeyLatency string
	LastCacheBypass     bool

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

func (m *FakeMetric) CacheHit(string) {
	m.LastCacheHit = true
}

func (m *FakeMetric) CacheMiss(string) {
	m.LastCacheMiss = true
}

func (m *FakeMetric) CacheInvalid(string) {
	m.LastCacheInvalid = true
}

func (m *FakeMetric) CacheError(key string, err string) {
	m.LastCacheKeyErr = key
	m.LastCacheErr = err
}

func (m *FakeMetric) CacheLatency(key string, duration time.Duration) {
	m.LastCacheKeyLatency = key
	m.LastCacheLatency = duration
}

func (m *FakeMetric) CacheBypassed() {
	m.LastCacheBypass = true
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
