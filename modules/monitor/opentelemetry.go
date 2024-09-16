package monitor

import (
	"context"
	"log"

	"github.com/gofiber/contrib/otelfiber/v2"
	"go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func (module *MonitorModule) initOpentelemetry() {
	module.tp = module.initTracer()
	module.app.Use(otelfiber.Middleware())
}

func (module *MonitorModule) destroyOpentelemetry() {
	if err := module.tp.Shutdown(context.Background()); err != nil {
		log.Printf("Error shutting down tracer provider: %v", err)
	}

	module.app.Use(otelfiber.Middleware())
}

func (module *MonitorModule) initTracer() *sdktrace.TracerProvider {
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpointURL(module.config.Getenv("OTLP_URL", "http://localhost:4318/v1/traces")),
	)
	if err != nil {
		log.Fatal(err)
	}
	appCode := module.config.Getenv("APP_CODE", "App")
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(appCode),
			)),
	)
	module.Service.setTracer(tp.Tracer(appCode))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}
