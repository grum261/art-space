package service

import (
	"art_space/internal/models"
	"context"
	"fmt"
)

type CommentsRepository interface {
	SelectCommentsByPost(ctx context.Context, postId int32) ([]models.Comment, error)
}

type Comment struct {
	repo CommentsRepository
}

func NewComment(repo CommentsRepository) *Comment {
	return &Comment{
		repo: repo,
	}
}

func (c *Comment) SelectCommentsByPost(ctx context.Context, postId int32) ([]models.Comment, error) {
	comments, err := c.repo.SelectCommentsByPost(ctx, postId)
	if err != nil {
		return nil, fmt.Errorf("реплозиторий комментариев - получение по id поста: %w", err)
	}

	return comments, nil
}
