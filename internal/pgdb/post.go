package pgdb

import (
	"art_space/internal/models"
	"context"

	"github.com/jackc/pgx/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Post репозиторий для взаимодействия с записями в базе
type Post struct {
	q *Queries
}

// NewPost создает новый инстанс Post репозитория
func NewPost(db *pgx.Conn) *Post {
	return &Post{
		q: New(db),
	}
}

// CreatePost вставляет новый пост в базу
func (p *Post) CreatePost(ctx context.Context, text string, authorId int) (int, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()

	postId, err := p.q.InsertPost(
		ctx, PostInsert{
			Text:     text,
			AuthorId: authorId,
		},
	)
	if err != nil {
		return 0, err
	}

	return postId, nil
}

// SelectAllPosts возвращает все посты из базы
func (p *Post) SelectAllPosts(ctx context.Context) ([]models.Post, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()

	posts, err := p.q.SelectAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	for i := range posts {
		comments, err := p.q.SelectCommentsOnPost(ctx, posts[i].Id)
		if err != nil {
			return nil, err
		}

		posts[i].Comments = comments
	}

	return posts, nil
}

// SelectPostsByAuthor возвращает все посты по автору
func (p *Post) SelectPostsByAuthor(ctx context.Context, authorName string) ([]models.Post, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()

	posts, err := p.q.SelectPostsByAuthor(ctx, authorName)
	if err != nil {
		return nil, err
	}

	for i := range posts {
		comments, err := p.q.SelectCommentsOnPost(ctx, posts[i].Id)
		if err != nil {
			return nil, err
		}

		posts[i].Comments = comments
	}

	return posts, nil
}

// SelectPostById возвращает пост по id
func (p *Post) SelectPostById(ctx context.Context, postId int) (models.Post, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()

	post, err := p.q.SelectPostById(ctx, postId)
	if err != nil {
		return models.Post{}, err
	}

	comments, err := p.q.SelectCommentsOnPost(ctx, postId)
	if err != nil {
		return models.Post{}, err
	}

	post.Comments = comments

	return post, nil
}

// UpdatePost обновляет пост по id
func (p *Post) UpdatePost(ctx context.Context, postId int, postText string) error {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()

	if err := p.q.UpdatePost(ctx, PostUpdate{
		Id:   postId,
		Text: postText,
	}); err != nil {
		return err
	}

	return nil
}

// DeletePost удаляет пост
func (p *Post) DeletePost(ctx context.Context, postId int) error {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()

	if err := p.q.DeletePost(ctx, postId); err != nil {
		return err
	}

	return nil
}
