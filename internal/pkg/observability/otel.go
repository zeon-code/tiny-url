package observability

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zeon-code/tiny-url/internal/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var serviceName string = "tiny_url"

type Observer interface {
	Logger() Logger
	Startup(ctx context.Context) error
	Shutdown(ctx context.Context) error
	Metric() (MetricClient, error)
	RegisterDB(dbStats metric.Registration)
}

type observer struct {
	AppVersion string
	Conf       config.Configuration

	tracer  *trace.TracerProvider
	metric  *sdkmetric.MeterProvider
	dbStats []metric.Registration
}

func NewObserver(appVersion string, conf config.Configuration) *observer {
	return &observer{
		AppVersion: appVersion,
		Conf:       conf,
	}
}

func (o *observer) Startup(ctx context.Context) error {
	conf := o.Conf.Metric()

	env, err := conf.Environment()

	if err != nil {
		return err
	}

	addr, err := conf.Host()

	if err != nil {
		return err
	}

	port, err := conf.Port()

	if err != nil {
		return err
	}

	exportTracerOptions := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%d", addr, port)),
	}

	exportMetricOptions := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(fmt.Sprintf("%s:%d", addr, port)),
	}

	if env == "local" {
		exportTracerOptions = append(exportTracerOptions, otlptracegrpc.WithInsecure())
		exportMetricOptions = append(exportMetricOptions, otlpmetricgrpc.WithInsecure())
	}

	tracerExporter, err := otlptracegrpc.New(ctx, exportTracerOptions...)

	if err != nil {
		return err
	}

	meterExporter, err := otlpmetricgrpc.New(ctx, exportMetricOptions...)

	if err != nil {
		return err
	}

	tracerResource, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.DeploymentEnvironment(env),
			semconv.ServiceVersion(o.AppVersion),
		),
	)

	if err != nil {
		return err
	}

	meterResource, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.DeploymentEnvironment(env),
			semconv.ServiceVersion(o.AppVersion),
		),
	)

	if err != nil {
		return err
	}

	o.tracer = trace.NewTracerProvider(
		trace.WithBatcher(tracerExporter, trace.WithBatchTimeout(1*time.Second)),
		trace.WithResource(tracerResource),
	)

	meterReader := sdkmetric.NewPeriodicReader(
		meterExporter,
		sdkmetric.WithInterval(1*time.Second),
	)

	o.metric = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(meterReader),
		sdkmetric.WithResource(meterResource),
	)

	otel.SetTracerProvider(o.tracer)
	otel.SetMeterProvider(o.metric)
	return nil
}

func (o *observer) Shutdown(ctx context.Context) error {
	var err error

	if o.tracer != nil {
		err = errors.Join(err, o.tracer.Shutdown(ctx))
	}

	if o.metric != nil {
		err = errors.Join(err, o.metric.Shutdown(ctx))
	}

	if o.dbStats != nil {
		for _, dbStats := range o.dbStats {
			err = errors.Join(err, dbStats.Unregister())
		}
	}

	return err
}

func (o *observer) Logger() Logger {
	return NewLogger(o.Conf.Log())
}

func (o *observer) Metric() (MetricClient, error) {
	return NewMetricClient(
		otel.Meter(
			serviceName,
			metric.WithInstrumentationVersion(o.AppVersion),
		),
	)
}

func (o *observer) RegisterDB(dbStats metric.Registration) {
	o.dbStats = append(o.dbStats, dbStats)
}
