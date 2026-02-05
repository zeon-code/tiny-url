package test

import "time"

type FakeDatadog struct {
	Err error

	LastIncrName string
	LastIncrTags []string
	LastIncrRate float64

	LastTimingName     string
	LastTimingDuration time.Duration
	LastTimingTags     []string
	LastTimingRate     float64
}

func NewFakeDatadog() *FakeDatadog {
	return &FakeDatadog{}
}

func (b *FakeDatadog) Incr(name string, tags []string, rate float64) error {
	b.LastIncrName = name
	b.LastIncrTags = tags
	b.LastIncrRate = rate

	return b.Err
}

func (b *FakeDatadog) Timing(name string, duration time.Duration, tags []string, rate float64) error {
	b.LastTimingName = name
	b.LastTimingDuration = duration
	b.LastTimingTags = tags
	b.LastTimingRate = rate

	return b.Err
}
