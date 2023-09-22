package tracing

import (
	"os"

	otelcontrib "go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func GetTracer() trace.Tracer {
	return otel.GetTracerProvider().Tracer(
		os.Getenv("TRACING_APP_NAME"),
		oteltrace.WithInstrumentationVersion(otelcontrib.Version()),
	)
}
