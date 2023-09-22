package tracing

var (
	TraceIdHeader = "X-B3-TraceId"
	SpanIdHeader  = "X-B3-SpanId"
)

//func CreateSpanContextFromTraceAndSpanId(traceId, spanId string) context.Context {
//	ctx := context.Background()
//
//	if traceId == "" || spanId == "" {
//		return ctx
//	}
//
//	traceID, err := trace.TraceIDFromHex(traceId)
//	if err != nil {
//		return ctx
//	}
//
//	spanID, err := trace.SpanIDFromHex(spanId)
//	if err != nil {
//		return ctx
//	}
//
//	cfg := trace.SpanContextConfig{
//		TraceID: traceID,
//		SpanID:  spanID,
//	}
//
//	spanContext := trace.NewSpanContext(cfg)
//	return trace.ContextWithRemoteSpanContext(ctx, spanContext)
//}
