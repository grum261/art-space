package rest

import "github.com/gofiber/fiber/v2"

type JsonResponse struct {
	Error  error       `json:"error"`
	Result interface{} `json:"result"`
}

func renderResponse(c *fiber.Ctx, result interface{}, statusCode int, err error) error {
	return c.Status(statusCode).JSON(JsonResponse{
		Error:  err,
		Result: result,
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
