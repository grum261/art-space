package rest

import (
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

func renderResponse(c *fiber.Ctx, result interface{}, statusCode int, err error) error {
	if err != nil {
		_, span := trace.SpanFromContext(c.Context()).TracerProvider().Tracer("").Start(c.Context(), "rest.renderErrResponse")
		span.RecordError(err)
		defer span.End()

		return c.Status(statusCode).JSON(map[string]interface{}{
			"error":  err.Error(),
			"result": nil,
		})
	}
	_, span := trace.SpanFromContext(c.Context()).TracerProvider().Tracer("").Start(c.Context(), "rest.renderResponseOK")
	defer span.End()

	return c.Status(statusCode).JSON(map[string]interface{}{
		"error":  nil,
		"result": result,
	})
}

func renderErrReponse(c *fiber.Ctx, err error) error {
	return renderResponse(c, nil, 500, err)
}

func renderResponseOK(c *fiber.Ctx, result interface{}) error {
	return renderResponse(c, result, 200, nil)
}

func renderResponseCreated(c *fiber.Ctx, result interface{}) error {
	return renderResponse(c, result, 201, nil)
}
