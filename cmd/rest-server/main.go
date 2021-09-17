package main

import (
	dockercontainer "art_space/internal/docker-container"
	"art_space/internal/envvar"
	"art_space/internal/models/service"
	"art_space/internal/pgdb"
	"art_space/internal/rest"
	"context"
	"log"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/export/metric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func main() {
	ctx := context.TODO()

	// TODO: сделать флаг для запуска без докера
	if err := dockercontainer.StartAllContainers(); err != nil {
		log.Fatal(err)
	}

	// TODO: добавить мигрирование

	db := pgdb.NewDB(ctx)
	defer db.Close(ctx)

	postRepo := pgdb.NewPost(db)
	svc := service.NewPost(postRepo)

	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 1,
		IdleTimeout:  time.Second * 1,
	})

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// TODO: нормальное логгирование для разных уровней
	app.Use(func(c *fiber.Ctx) error {
		logger.Info(c.Method(), zap.Time("time", time.Now()), zap.String("url", c.BaseURL()))

		return c.Next()
	})

	rest.RegisterOpenApi(app)

	app.Static("/", "./assets/swagger-ui")

	rest.NewPostHandler(svc).RegisterRoutes(app)

	if err := initTracer(); err != nil {
		log.Fatal(err)
	}

	app.Use(adaptor.HTTPMiddleware(otelmux.Middleware("artspace-api-server")))

	log.Fatal(app.Listen(":8000"))
}

func initTracer() error {
	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
		return err
	}

	promExport, err := prometheus.New(prometheus.Config{}, controller.New(
		processor.New(
			simple.NewWithHistogramDistribution(),
			metric.CumulativeExportKindSelector(),
			processor.WithMemory(true),
		),
	))
	if err != nil {
		return err
	}

	global.SetMeterProvider(promExport.MeterProvider())

	jaegerEndpoint, _ := envvar.VaultProvider.Get("JAEGER_ENDPOINT")

	jaegerExporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		return err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSyncer(jaegerExporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL, attribute.KeyValue{
				Key:   semconv.ServiceNameKey,
				Value: attribute.StringValue("artspace-api-server"),
			},
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return nil
}
