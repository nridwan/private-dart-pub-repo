package monitor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type MonitorSpan = trace.Span
type MonitorTracer = trace.Tracer

type MonitorService interface {
	StartTraceSpan(context context.Context, title string, attributes map[string]interface{}) (context.Context, MonitorSpan)
	SetCurrentSpanAttributes(context context.Context, attributes map[string]interface{})
	setTracer(tracer MonitorTracer)
}

type monitorServiceImpl struct {
	tracer MonitorTracer
}

func provideMonitorService() MonitorService {
	return &monitorServiceImpl{}
}

func mapToAttributes(attributes map[string]interface{}) []attribute.KeyValue {
	var result []attribute.KeyValue
	for key, value := range attributes {
		switch v := value.(type) {
		case string:
			result = append(result, attribute.String(key, v))
		case int:
			result = append(result, attribute.Int64(key, int64(v)))
		case float64:
			result = append(result, attribute.Float64(key, v))
		case bool:
			result = append(result, attribute.Bool(key, v))
		default:
			// Handle other types as needed
			result = append(result, attribute.String(key, fmt.Sprintf("%v", value)))
		}
	}
	return result
}

func (service *monitorServiceImpl) StartTraceSpan(context context.Context, title string, attributes map[string]interface{}) (context.Context, MonitorSpan) {
	return service.tracer.Start(context, title, trace.WithAttributes(mapToAttributes(attributes)...))
}

func (service *monitorServiceImpl) SetCurrentSpanAttributes(context context.Context, attributes map[string]interface{}) {
	trace.SpanFromContext(context).SetAttributes(mapToAttributes(attributes)...)
}

func (service *monitorServiceImpl) setTracer(tracer MonitorTracer) {
	service.tracer = tracer
}
