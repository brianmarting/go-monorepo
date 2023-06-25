package tracing

import (
	otelcontrib "go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	name = "trace/go-backend"
)

func GetTracer() trace.Tracer {
	return otel.GetTracerProvider().Tracer(
		name,
		oteltrace.WithInstrumentationVersion(otelcontrib.SemVersion()),
	)
}
