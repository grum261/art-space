package rest

import (
	"art_space/internal/models"
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

//go:generate counterfeiter -o resttesting/post_service.gen.go . PostService

type PostService interface {
	CreatePost(ctx context.Context, text string, authorId int) (int, error)
	UpdatePost(ctx context.Context, id int, text string) error
	DeletePost(ctx context.Context, id int) error
	SelectPostById(ctx context.Context, id int) (models.Post, error)
	SelectAllPosts(ctx context.Context) ([]models.Post, error)
}

type PostHandler struct {
	svc PostService
}

func NewPostHandler(svc PostService) *PostHandler {
	return &PostHandler{
		svc: svc,
	}
}

type Post struct {
	Id       int       `json:"id"`
	Text     string    `json:"text"`
	Dates    Dates     `json:"dates"`
	Comments []Comment `json:"comments,omitempty"`
	Author   Author    `json:"author"`
}

type Author struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type CreateUpdatePostRequest struct {
	Text     string `json:"text"`
	AuthorId int    `json:"author"`
}

type Comment struct {
	Id     int    `json:"id"`
	Text   string `json:"text"`
	Dates  Dates  `json:"dates"`
	Author Author `json:"author"`
}

type Dates struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (ph *PostHandler) RegisterRoutes(app *fiber.App) {
	app.Post("/posts", ph.createHandler)
	app.Get("/posts", ph.selectAllHandler)
	app.Get("/posts/:postId", ph.selectByIdHandler)
	app.Put("/posts/:postId", ph.updateHandler)
	app.Delete("/posts/:postId", ph.deleteHandler)
}

func (ph *PostHandler) createHandler(c *fiber.Ctx) error {
	req := &CreateUpdatePostRequest{}

	if err := c.BodyParser(req); err != nil {
		return renderErrReponse(c, err)
	}

	postId, err := ph.svc.CreatePost(c.Context(), req.Text, req.AuthorId)
	if err != nil {
		return renderErrReponse(c, err)
	}

	return renderResponseCreated(c, postId)
}

func (ph *PostHandler) selectAllHandler(c *fiber.Ctx) error {
	posts, err := ph.svc.SelectAllPosts(c.Context())
	if err != nil {
		return renderErrReponse(c, err)
	}

	response := []*Post{}

	for _, post := range posts {
		p := &Post{
			Id:   post.Id,
			Text: post.Text,
			Dates: Dates{
				CreatedAt: post.Dates.CreatedAt,
				UpdatedAt: post.Dates.UpdatedAt,
			},
			Author: Author{
				Id:     post.Author.Id,
				Name:   post.Author.Name,
				Avatar: post.Author.Avatar,
			},
		}

		for i, comment := range post.Comments {
			p.Comments[i] = Comment{
				Id:   comment.Id,
				Text: comment.Text,
				Dates: Dates{
					CreatedAt: comment.Dates.CreatedAt,
					UpdatedAt: comment.Dates.UpdatedAt,
				},
				Author: Author{
					Id:     comment.Author.Id,
					Name:   comment.Author.Name,
					Avatar: comment.Author.Avatar,
				},
			}
		}

		response = append(response, p)
	}

	return renderResponseOK(c, response)
}

func (ph *PostHandler) selectByIdHandler(c *fiber.Ctx) error {
	postId, err := strconv.Atoi(c.Params("postId"))
	if err != nil {
		return renderErrReponse(c, err)
	}

	post, err := ph.svc.SelectPostById(c.Context(), postId)
	if err != nil {
		renderErrReponse(c, err)
	}

	response := &Post{
		Id:   post.Id,
		Text: post.Text,
		Dates: Dates{
			CreatedAt: post.Dates.CreatedAt,
			UpdatedAt: post.Dates.UpdatedAt,
		},
		Author: Author{
			Id:     post.Author.Id,
			Name:   post.Author.Name,
			Avatar: post.Author.Avatar,
		},
	}

	for i, comment := range post.Comments {
		response.Comments[i] = Comment{
			Id:   comment.Id,
			Text: comment.Text,
			Dates: Dates{
				CreatedAt: comment.Dates.CreatedAt,
				UpdatedAt: comment.Dates.UpdatedAt,
			},
			Author: Author{
				Id:     comment.Author.Id,
				Name:   comment.Author.Name,
				Avatar: comment.Author.Avatar,
			},
		}
	}

	return renderResponseOK(c, response)
}

func (ph *PostHandler) updateHandler(c *fiber.Ctx) error {
	postId, err := strconv.Atoi(c.Params("postId"))
	if err != nil {
		return renderErrReponse(c, err)
	}

	req := &CreateUpdatePostRequest{}

	if err := c.BodyParser(req); err != nil {
		return renderErrReponse(c, err)
	}

	if err := ph.svc.UpdatePost(c.Context(), postId, req.Text); err != nil {
		return renderResponseOK(c, err)
	}

	return renderResponseOK(c, postId)
}

func (ph *PostHandler) deleteHandler(c *fiber.Ctx) error {
	postId, err := strconv.Atoi(c.Params("postId"))
	if err != nil {
		return renderErrReponse(c, err)
	}

	if err := ph.svc.DeletePost(c.Context(), postId); err != nil {
		return renderErrReponse(c, err)
	}

	return renderResponseOK(c, postId)
}
