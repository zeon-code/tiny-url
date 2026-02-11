package observability

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/zeon-code/tiny-url/internal/pkg/config"
	"go.opentelemetry.io/otel/trace"
)

type Logger interface {
	With(...any) Logger

	Debug(context.Context, string, ...any)
	Info(context.Context, string, ...any)
	Warn(context.Context, string, ...any)
	Error(context.Context, string, ...any)
}

func NewLogger(conf config.Log) Logger {
	var level slog.Level

	switch strings.ToLower(conf.Level()) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	return NewOtelLogger(level)
}

type OtelLogger struct {
	logger *slog.Logger
}

func NewOtelLogger(level slog.Level) OtelLogger {
	return OtelLogger{
		logger: slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: level,
				},
			),
		),
	}
}

func (l OtelLogger) Debug(ctx context.Context, msg string, args ...any) {
	l.withTrace(ctx).Debug(msg, args...)
}

func (l OtelLogger) Info(ctx context.Context, msg string, args ...any) {
	l.withTrace(ctx).Info(msg, args...)
}

func (l OtelLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.withTrace(ctx).Warn(msg, args...)
}

func (l OtelLogger) Error(ctx context.Context, msg string, args ...any) {
	l.withTrace(ctx).Error(msg, args...)
}

func (l OtelLogger) With(args ...any) Logger {
	return OtelLogger{
		logger: l.logger.With(args...),
	}
}

func (l OtelLogger) withTrace(ctx context.Context) *slog.Logger {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return l.logger
	}

	sc := span.SpanContext()

	return l.logger.With(
		slog.String("trace_id", sc.TraceID().String()),
		slog.String("span_id", sc.SpanID().String()),
	)
}
