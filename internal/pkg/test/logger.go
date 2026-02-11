package test

import (
	"context"

	"github.com/zeon-code/tiny-url/internal/pkg/observability"
)

type FakeLogger struct{}

func (l FakeLogger) With(args ...any) observability.Logger              { return FakeLogger{} }
func (l FakeLogger) Debug(ctx context.Context, msg string, args ...any) {}
func (l FakeLogger) Info(ctx context.Context, msg string, args ...any)  {}
func (l FakeLogger) Warn(ctx context.Context, msg string, args ...any)  {}
func (l FakeLogger) Error(ctx context.Context, msg string, args ...any) {}
