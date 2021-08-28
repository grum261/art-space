package pgdb

import (
	"context"
	"fmt"
)

// TODO: возможно сделать, чтобы выражаения подготавливались перед выполнением запроса, а не при инициализации базы

// initPostRequests подготавливает все запросы к базе для постов
func initPostRequests(ctx context.Context, db PGDB) error {
	if _, err := db.Prepare(
		ctx, "Insert.Post.",
		`INSERT INTO posts (text, author_id, created_at) VALUES ($1, $2, current_timestamp) RETURNING id`,
	); err != nil {
		return fmt.Errorf("ошибка подготовки запроса Insert.Post.: %w", err)
	}

	if _, err := db.Prepare(ctx, "Update.Post.Id", `UPDATE posts SET text = $2 WHERE id = $1`); err != nil {
		return fmt.Errorf("ошибка подготовки запроса Update.Post.Id: %w", err)
	}

	if _, err := db.Prepare(
		ctx, "Select.Post.All",
		`SELECT p.id, p.text, p.created_at, p.updated_at, u.id as author_id, u.username as author_name, u.avatar
		FROM posts p
		INNER JOIN users u ON p.author_id = u.id`,
	); err != nil {
		return fmt.Errorf("ошибка подготовки запроса Select.Post.All: %w", err)
	}

	if _, err := db.Prepare(
		ctx, "Select.Post.Id",
		`SELECT p.id, p.text, p.created_at, p.updated_at, u.id as author_id, u.username as author_name, u.avatar
		FROM posts p
		INNER JOIN users u ON p.author_Id = u.id`,
	); err != nil {
		return fmt.Errorf("ошибка подготовки запроса Select.Post.Id: %w", err)
	}

	if _, err := db.Prepare(
		ctx, "Select.Post.Author",
		`SELECT p.id, p.text, p.created_at, p.updated_at, u.id as author_id, u.username as author_name, u.avatar
		FROM posts p
		INNER JOIN users u ON p.author_Id = u.id
		WHERE u.username = $1
		ORDER BY p.created_at`,
	); err != nil {
		return fmt.Errorf("ошибка подготовки запроса Select.Post.Author: %w", err)
	}

	if _, err := db.Prepare(ctx, "Delete.Post.Id", `DELETE FROM posts WHERE id = $1`); err != nil {
		return fmt.Errorf("ошибка подготовки запроса Delete.Post.Id: %w", err)
	}

	return nil
}

// initCommentRequests подготавливает все запросы к базе для комментариев
func initCommentRequests(ctx context.Context, db PGDB) error {
	if _, err := db.Prepare(
		ctx, "Select.Comment.Post",
		`SELECT c.id, c.text, c.created_at, c.updated_at, u.id as author_id, u.username as author_name, u.avatar
		FROM comments c
		INNER JOIN users u ON c.author_id = u.id
		INNER JOIN posts p ON c.post_id = p.id
		WHERE c.post_id = $1`,
	); err != nil {
		return fmt.Errorf("ошибка подготовки запроса Select.Comment.Post: %w", err)
	}

	return nil
}
