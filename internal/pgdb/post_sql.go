package pgdb

import (
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
	row := q.db.QueryRow(ctx, "Insert.Post.", _in.Text, _in.AuthorId)

	var postId int

	if err := row.Scan(&postId); err != nil {
		return 0, fmt.Errorf("(InsertPost) ошибка создания поста: %w", err)
	}

	return postId, nil
}

func (q *Queries) SelectAllPosts(ctx context.Context) ([]PostReturn, error) {
	rows, err := q.db.Query(ctx, "Select.Post.All")
	if err != nil {
		return nil, fmt.Errorf("(SelectAllPosts) ошибка получения всех постов: %w", err)
	}
	defer rows.Close()

	var _out []PostReturn

	for rows.Next() {
		var r PostReturn

		if err := rows.Scan(
			&r.Id, &r.Text, &r.CreatedAt, &r.UpdatedAt,
			&r.AuthorId, &r.AuthorName, &r.AuthorAvatar,
		); err != nil {
			return nil, fmt.Errorf("(SelectAllPosts) ошибка сканирования строк: %w", err)
		}

		_out = append(_out, r)
	}

	return _out, nil
}

func (q *Queries) SelectPostById(ctx context.Context, id int) (PostReturn, error) {
	row := q.db.QueryRow(ctx, "Select.Post.Id", id)

	var r PostReturn

	if err := row.Scan(
		&r.Id, &r.Text, &r.CreatedAt, &r.UpdatedAt,
		&r.AuthorId, &r.AuthorName, &r.AuthorAvatar,
	); err != nil {
		return PostReturn{}, fmt.Errorf("(SelectPostById) ошибка сканирования строк: %w", err)
	}

	return r, nil
}

func (q *Queries) SelectPostsByAuthor(ctx context.Context, username string) ([]PostReturn, error) {
	rows, err := q.db.Query(ctx, "Select.Post.Author", username)
	if err != nil {
		return nil, fmt.Errorf("(SelectPostsByAuthor) ошибка получения всех постов по автору: %w", err)
	}
	defer rows.Close()

	var _out []PostReturn

	for rows.Next() {
		var r PostReturn

		if err := rows.Scan(
			&r.Id, &r.Text, &r.CreatedAt, &r.UpdatedAt,
			&r.AuthorId, &r.AuthorName, &r.AuthorAvatar,
		); err != nil {
			return nil, err
		}

		_out = append(_out, r)
	}

	return _out, nil
}

func (q *Queries) DeletePost(ctx context.Context, id int) error {
	if _, err := q.db.Exec(ctx, "Delete.Post.Id", id); err != nil {
		return fmt.Errorf("(DeletePost) ошибка удаления поста: %w", err)
	}

	return nil
}

func (q *Queries) UpdatePost(ctx context.Context, _in PostUpdate) error {
	if _, err := q.db.Exec(ctx, "Update.Post.Id", _in.Id, _in.Text); err != nil {
		return fmt.Errorf("(UpdatePost) ошибка обновления поста: %w", err)
	}

	return nil
}
