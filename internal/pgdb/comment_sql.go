package pgdb

import (
	"art_space/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgtype"
)

type CommentReturn struct {
	Id           int
	Text         string
	CreatedAt    pgtype.Timestamp
	UpdatedAt    pgtype.Timestamp
	AuthorId     int
	AuthorName   string
	AuthorAvatar string
}

func (q *Queries) SelectCommentsOnPost(ctx context.Context, postId int) ([]models.Comment, error) {
	rows, err := q.db.Query(ctx, selectCommentsByPost, postId)
	if err != nil {
		return nil, fmt.Errorf("(pgdb.Queries.SelectCommentsOnPost) ошибка получение комментариев к посту: %w", err)
	}
	defer rows.Close()

	var _out []models.Comment

	for rows.Next() {
		var r CommentReturn

		if err := rows.Scan(
			&r.Id, &r.Text, &r.CreatedAt, &r.UpdatedAt,
			&r.AuthorId, &r.AuthorName, &r.AuthorAvatar,
		); err != nil {
			return nil, fmt.Errorf("(pgdb.Queries.SelectCommentsOnPost) ошибка сканирования строк: %w", err)
		}

		_out = append(
			_out, models.Comment{
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
