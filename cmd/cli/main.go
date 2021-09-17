package main

import (
	"art_space/internal/envvar"
	"art_space/pkg/openapi3"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	ctx := context.Background()

	initTracer()

	clientOA := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	client, err := openapi3.NewClientWithResponses("http://172.25.122.55:8000", openapi3.WithHTTPClient(&clientOA))
	if err != nil {
		log.Fatalf("не удалось запустить клиент swagger-ui: %s", err)
	}

	newPtrStr := func(s string) *string {
		return &s
	}

	newPtrInt32 := func(i int32) *int32 {
		return &i
	}

	respCreate, err := client.CreatePostWithResponse(ctx, openapi3.CreatePostJSONRequestBody{
		AuthorId: newPtrInt32(1),
		Text:     newPtrStr("dksopakds[aokd"),
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.UpdatePostWithResponse(ctx, *respCreate.JSON201.Result, openapi3.UpdatePostJSONRequestBody{
		AuthorId: newPtrInt32(2),
		Text:     newPtrStr("yoyoyoyoyo"),
	})
	if err != nil {
		log.Fatal(err)
	}

	respSelect, err := client.SelectAllWithResponse(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*respSelect.JSON200.Result)

	time.Sleep(10 * time.Second)
}

func initTracer() {
	jaegerEndpoint, _ := envvar.VaultProvider.Get("JAEGER_ENDPOINT")

	jaegerExporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSyncer(jaegerExporter),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}
