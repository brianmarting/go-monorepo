package tracing

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func InitTracerProvider() *trace.TracerProvider {
	ctx := context.Background()
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(os.Getenv("OTLP_ENDPOINT")),
		otlptracehttp.WithInsecure(),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create trace exporter")
	}

	batchSpanProcessor := trace.NewBatchSpanProcessor(exporter)

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", "go-backend"),
		),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create otel resource")
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithSpanProcessor(batchSpanProcessor),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tracerProvider
}
