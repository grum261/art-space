package pgdb

import (
	"art_space/internal/models"
	"context"

	"github.com/jackc/pgx/v4"
)

type Post struct {
	q *Queries
}

func NewPost(db *pgx.Conn) *Post {
	return &Post{
		q: New(db),
	}
}

func (p *Post) CreatePost(ctx context.Context, text string, authorId int) (int, error) {
	postId, err := p.q.InsertPost(ctx, PostInsert{
		Text:     text,
		AuthorId: authorId,
	})
	if err != nil {
		return 0, err
	}

	return postId, nil
}

func (p *Post) SelectAllPosts(ctx context.Context) ([]models.Post, error) {
	posts, err := p.q.SelectAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	var _out []models.Post

	for _, post := range posts {
		pm := models.Post{
			Id:   post.Id,
			Text: post.Text,
			Dates: models.Dates{
				CreatedAt: post.CreatedAt.Time,
				UpdatedAt: post.UpdatedAt.Time,
			},
			Author: models.Author{
				Id:     post.AuthorId,
				Name:   post.AuthorName,
				Avatar: post.AuthorAvatar,
			},
		}

		comments, err := p.q.SelectCommentsOnPost(ctx, post.Id)
		if err != nil {
			return nil, err
		}

		for _, comment := range comments {
			pm.Comments = append(pm.Comments, models.Comment{
				Id:   comment.Id,
				Text: comment.Text,
				Dates: models.Dates{
					CreatedAt: comment.CreatedAt.Time,
					UpdatedAt: comment.UpdatedAt.Time,
				},
				Author: models.Author{
					Id:     comment.AuthorId,
					Name:   comment.AuthorName,
					Avatar: comment.AuthorAvatar,
				},
			})
		}

		_out = append(_out, pm)
	}

	return _out, nil
}

func (p *Post) SelectPostsByAuthor(ctx context.Context, authorName string) ([]models.Post, error) {
	posts, err := p.q.SelectPostsByAuthor(ctx, authorName)
	if err != nil {
		return nil, err
	}

	var _out []models.Post

	for _, post := range posts {
		pm := models.Post{
			Id:   post.Id,
			Text: post.Text,
			Dates: models.Dates{
				CreatedAt: post.CreatedAt.Time,
				UpdatedAt: post.UpdatedAt.Time,
			},
			Author: models.Author{
				Id:     post.AuthorId,
				Name:   authorName,
				Avatar: post.AuthorAvatar,
			},
		}

		comments, err := p.q.SelectCommentsOnPost(ctx, post.Id)
		if err != nil {
			return nil, err
		}

		for _, comment := range comments {
			pm.Comments = append(pm.Comments, models.Comment{
				Id:   comment.Id,
				Text: comment.Text,
				Dates: models.Dates{
					CreatedAt: comment.CreatedAt.Time,
					UpdatedAt: comment.UpdatedAt.Time,
				},
				Author: models.Author{
					Id:     comment.AuthorId,
					Name:   comment.AuthorName,
					Avatar: comment.AuthorAvatar,
				},
			})
		}

		_out = append(_out, pm)
	}

	return _out, nil
}

func (p *Post) SelectPostById(ctx context.Context, postId int) (models.Post, error) {
	post, err := p.q.SelectPostById(ctx, postId)
	if err != nil {
		return models.Post{}, err
	}

	pm := models.Post{
		Id:   postId,
		Text: post.Text,
		Dates: models.Dates{
			CreatedAt: post.CreatedAt.Time,
			UpdatedAt: post.UpdatedAt.Time,
		},
		Author: models.Author{
			Id:     post.AuthorId,
			Name:   post.AuthorName,
			Avatar: post.AuthorAvatar,
		},
	}

	comments, err := p.q.SelectCommentsOnPost(ctx, postId)
	if err != nil {
		return models.Post{}, err
	}

	for _, comment := range comments {
		pm.Comments = append(pm.Comments, models.Comment{
			Id:   comment.Id,
			Text: comment.Text,
			Dates: models.Dates{
				CreatedAt: comment.CreatedAt.Time,
				UpdatedAt: comment.UpdatedAt.Time,
			},
			Author: models.Author{
				Id:     comment.AuthorId,
				Name:   comment.AuthorName,
				Avatar: comment.AuthorAvatar,
			},
		})
	}

	return pm, nil
}

func (p *Post) UpdatePost(ctx context.Context, postId int, postText string) error {
	if err := p.q.UpdatePost(ctx, PostUpdate{
		Id:   postId,
		Text: postText,
	}); err != nil {
		return err
	}

	return nil
}

func (p *Post) DeletePost(ctx context.Context, postId int) error {
	if err := p.q.DeletePost(ctx, postId); err != nil {
		return err
	}

	return nil
}
