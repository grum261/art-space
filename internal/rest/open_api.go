package rest

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	"github.com/gofiber/fiber/v2"
)

//go:generate go run ../../cmd/openapi-gen/main.go -path .
//go:generate oapi-codegen -package openapi3 -generate types -o ../../pkg/openapi3/task_types.gen.go openapi3.yaml
//go:generate oapi-codegen -package openapi3 -generate client -o ../../pkg/openapi3/client.gen.go     openapi3.yaml

func NewOpenApi3() openapi3.T {
	swagger := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "Art space API",
			Description: "REST API used for interacting with Art space app",
		},
		Servers: openapi3.Servers{
			{
				URL:         "192.168.183.84:8000",
				Description: "Local dev",
			},
		},
	}

	swagger.Components.Schemas = openapi3.Schemas{
		"Author": openapi3.NewSchemaRef(
			"", openapi3.NewObjectSchema().WithProperties(
				map[string]*openapi3.Schema{
					"id":       openapi3.NewInt32Schema(),
					"username": openapi3.NewStringSchema(),
					"avatar":   openapi3.NewStringSchema().WithNullable(),
				},
			),
		),
		"Dates": openapi3.NewSchemaRef(
			"", openapi3.NewObjectSchema().WithProperties(
				map[string]*openapi3.Schema{
					"createdAt": openapi3.NewDateTimeSchema(),
					"updatedAt": openapi3.NewDateTimeSchema(),
				},
			),
		),
		"Comments": openapi3.NewSchemaRef(
			"", openapi3.NewArraySchema().WithProperties(
				map[string]*openapi3.Schema{
					"id":   openapi3.NewInt32Schema(),
					"text": openapi3.NewStringSchema(),
				},
			).WithPropertyRef("dates", &openapi3.SchemaRef{Ref: "#/components/schemas/Dates"}).
				WithPropertyRef("author", &openapi3.SchemaRef{Ref: "#/components/schemas/Author"}),
		),
		"Post": openapi3.NewSchemaRef(
			"", openapi3.NewObjectSchema().WithProperties(
				map[string]*openapi3.Schema{
					"id":   openapi3.NewInt32Schema(),
					"text": openapi3.NewStringSchema(),
				},
			).WithPropertyRef("comments", &openapi3.SchemaRef{Ref: "#/components/schemas/Comments"}).
				WithPropertyRef("dates", &openapi3.SchemaRef{Ref: "#/components/schemas/Dates"}).
				WithPropertyRef("author", &openapi3.SchemaRef{Ref: "#/components/schemas/Author"}),
		),
	}

	swagger.Components.RequestBodies = openapi3.RequestBodies{
		"CreateUpdatePostRequest": &openapi3.RequestBodyRef{
			Value: openapi3.NewRequestBody().
				WithDescription("Тело запроса для создания и обновления поста").
				WithRequired(true).
				WithJSONSchema(openapi3.NewSchema().WithProperties(
					map[string]*openapi3.Schema{
						"text":     openapi3.NewStringSchema().WithMinLength(1),
						"authorId": openapi3.NewInt32Schema(),
					},
				)),
		},
	}

	swagger.Components.Responses = openapi3.Responses{
		"ErrorResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Ответ при ошибке").
				WithContent(openapi3.NewContentWithJSONSchema(
					openapi3.NewSchema().WithProperties(
						map[string]*openapi3.Schema{
							"error":  openapi3.NewStringSchema(),
							"result": openapi3.NewInt32Schema(),
						},
					),
				)),
		},
		"CreateUpdatePostResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Ответ при успешном создании или обновлении поста").
				WithContent(openapi3.NewContentWithJSONSchema(
					openapi3.NewSchema().WithProperties(
						map[string]*openapi3.Schema{
							"error":  openapi3.NewStringSchema(),
							"result": openapi3.NewInt32Schema(),
						},
					),
				)),
		},
		"SelectPostByIdResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Ответ при успешном получении поста по id").
				WithContent(openapi3.NewContentWithJSONSchema(
					openapi3.NewSchema().
						WithProperty("error", openapi3.NewStringSchema()).
						WithPropertyRef("result", &openapi3.SchemaRef{Ref: "#/components/schemas/Post"}),
				)),
		},
		"SelectAllPostsResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Ответ при успешном получении всех постов").
				WithContent(openapi3.NewContentWithJSONSchema(
					&openapi3.Schema{
						Properties: openapi3.Schemas{
							"error":  &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
							"result": openapi3.NewSchemaRef("", &openapi3.Schema{Items: &openapi3.SchemaRef{Ref: "#/components/schemas/Post"}}),
						},
					},
				)),
		},
	}

	swagger.Paths = openapi3.Paths{
		"/posts": &openapi3.PathItem{
			Post: &openapi3.Operation{
				OperationID: "CreatePost",
				RequestBody: &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/CreateUpdatePostRequest",
				},
				Responses: openapi3.Responses{
					"400": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"201": &openapi3.ResponseRef{
						Ref: "#/components/responses/CreateUpdatePostResponse",
					},
				},
			},
			Get: &openapi3.Operation{
				OperationID: "SelectAll",
				Responses: openapi3.Responses{
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/SelectAllPostsResponse",
					},
				},
			},
		},
		"/posts/{postId}": &openapi3.PathItem{
			Get: &openapi3.Operation{
				OperationID: "SelectById",
				Parameters: []*openapi3.ParameterRef{
					{
						Value: openapi3.NewPathParameter("postId").
							WithSchema(openapi3.NewInt32Schema()),
					},
				},
				Responses: openapi3.Responses{
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/SelectPostByIdResponse",
					},
				},
			},
			Put: &openapi3.Operation{
				OperationID: "UpdatePost",
				Parameters: []*openapi3.ParameterRef{
					{
						Value: openapi3.NewPathParameter("postId").
							WithSchema(openapi3.NewInt32Schema()),
					},
				},
				RequestBody: &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/CreateUpdatePostRequest",
				},
				Responses: openapi3.Responses{
					"400": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/CreateUpdatePostResponse",
					},
				},
			},
			Delete: &openapi3.Operation{
				OperationID: "DeletePost",
				Parameters: []*openapi3.ParameterRef{
					{
						Value: openapi3.NewPathParameter("postId").
							WithSchema(openapi3.NewInt32Schema()),
					},
				},
				Responses: openapi3.Responses{
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/CreateUpdatePostResponse",
					},
				},
			},
		},
	}

	return swagger
}

func RegisterOpenApi(app *fiber.App) {
	swagger := NewOpenApi3()

	app.Get("/openapi3.json", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(&swagger)
	})

	app.Get("/openapi3.yaml", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/x-yml")

		data, err := yaml.Marshal(&swagger)
		if err != nil {
			return err
		}

		_, err = c.Status(200).Write(data)

		return err
	})
}
