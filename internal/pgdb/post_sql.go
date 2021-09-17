package pgdb

import (
	"art_space/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgtype"
)

type PostInsert struct {
	Text     string
	AuthorId int
}

type PostUpdate struct {
	Id   int
	Text string
}

type PostReturn struct {
	Id           int
	Text         string
	CreatedAt    pgtype.Timestamp
	UpdatedAt    pgtype.Timestamp
	AuthorId     int
	AuthorName   string
	AuthorAvatar string
}

func (q *Queries) InsertPost(ctx context.Context, _in PostInsert) (int, error) {
	row := q.db.QueryRow(ctx, insertPost, _in.Text, _in.AuthorId)

	var postId int

	if err := row.Scan(&postId); err != nil {
		return 0, fmt.Errorf("(pgdb.Queries.InsertPost) ошибка создания поста: %w", err)
	}

	return postId, nil
}

func (q *Queries) SelectAllPosts(ctx context.Context) ([]models.Post, error) {
	rows, err := q.db.Query(ctx, selectAllPosts)
	if err != nil {
		return nil, fmt.Errorf("(pgdb.Queries.SelectAllPosts) ошибка получения всех постов: %w", err)
	}
	defer rows.Close()

	var _out []models.Post

	for rows.Next() {
		var r PostReturn

		if err := rows.Scan(
			&r.Id, &r.Text, &r.CreatedAt, &r.UpdatedAt,
			&r.AuthorId, &r.AuthorName, &r.AuthorAvatar,
		); err != nil {
			return nil, fmt.Errorf("(pgdb.Queries.SelectAllPosts) ошибка сканирования строк: %w", err)
		}

		_out = append(
			_out, models.Post{
				Id:   r.Id,
				Text: r.Text,
				Dates: models.Dates{
					CreatedAt: r.CreatedAt.Time,
					UpdatedAt: r.UpdatedAt.Time,
				},
				Author: models.Author{
					Id:     r.AuthorId,
					Name:   r.AuthorName,
					Avatar: r.AuthorAvatar,
				},
			},
		)
	}

	return _out, nil
}

func (q *Queries) SelectPostById(ctx context.Context, id int) (models.Post, error) {
	row := q.db.QueryRow(ctx, selectPostById, id)

	var r PostReturn

	if err := row.Scan(
		&r.Id, &r.Text, &r.CreatedAt, &r.UpdatedAt,
		&r.AuthorId, &r.AuthorName, &r.AuthorAvatar,
	); err != nil {
		return models.Post{}, fmt.Errorf("(pgdb.Queries.SelectPostById) ошибка сканирования строк: %w", err)
	}

	return models.Post{
		Id:   r.Id,
		Text: r.Text,
		Dates: models.Dates{
			CreatedAt: r.CreatedAt.Time,
			UpdatedAt: r.UpdatedAt.Time,
		},
		Author: models.Author{
			Id:     r.AuthorId,
			Name:   r.AuthorName,
			Avatar: r.AuthorAvatar,
		},
	}, nil
}

func (q *Queries) SelectPostsByAuthor(ctx context.Context, username string) ([]models.Post, error) {
	rows, err := q.db.Query(ctx, selectPostByAuthor, username)
	if err != nil {
		return nil, fmt.Errorf("(pgdb.Queries.SelectPostsByAuthor) ошибка получения всех постов по автору: %w", err)
	}
	defer rows.Close()

	var _out []models.Post

	for rows.Next() {
		var r PostReturn

		if err := rows.Scan(
			&r.Id, &r.Text, &r.CreatedAt, &r.UpdatedAt,
			&r.AuthorId, &r.AuthorName, &r.AuthorAvatar,
		); err != nil {
			return nil, err
		}

		_out = append(
			_out, models.Post{
				Id:   r.Id,
				Text: r.Text,
				Dates: models.Dates{
					CreatedAt: r.CreatedAt.Time,
					UpdatedAt: r.UpdatedAt.Time,
				},
				Author: models.Author{
					Id:     r.AuthorId,
					Name:   r.AuthorName,
					Avatar: r.AuthorAvatar,
				},
			},
		)
	}

	return _out, nil
}

func (q *Queries) DeletePost(ctx context.Context, id int) error {
	if _, err := q.db.Exec(ctx, deletePostById, id); err != nil {
		return fmt.Errorf("(pgdb.Queries.DeletePost) ошибка удаления поста: %w", err)
	}

	return nil
}

func (q *Queries) UpdatePost(ctx context.Context, _in PostUpdate) error {
	if _, err := q.db.Exec(ctx, updatePostById, _in.Id, _in.Text); err != nil {
		return fmt.Errorf("(pgdb.Queries.UpdatePost) ошибка обновления поста: %w", err)
	}

	return nil
}
