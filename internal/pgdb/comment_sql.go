package pgdb

import (
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

func (q *Queries) SelectCommentsOnPost(ctx context.Context, postId int) ([]CommentReturn, error) {
	rows, err := q.db.Query(ctx, "Select.Comment.Post", postId)
	if err != nil {
		return nil, fmt.Errorf("(SelectCommentsOnPost) ошибка получение комментариев к посту: %w", err)
	}
	defer rows.Close()

	var _out []CommentReturn

	for rows.Next() {
		var r CommentReturn

		if err := rows.Scan(
			&r.Id, &r.Text, &r.CreatedAt, &r.UpdatedAt,
			&r.AuthorId, &r.AuthorName, &r.AuthorAvatar,
		); err != nil {
			return nil, fmt.Errorf("(SelectCommentsOnPost) ошибка сканирования строк: %w", err)
		}

		_out = append(_out, r)
	}

	return _out, nil
}
