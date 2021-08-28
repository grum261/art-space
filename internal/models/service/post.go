package service

import (
	"art_space/internal/models"
	"context"
	"fmt"
)

type PostRepository interface {
	CreatePost(ctx context.Context, text string, authorId int) (int, error)
	UpdatePost(ctx context.Context, id int, text string) error
	DeletePost(ctx context.Context, id int) error
	SelectPostById(ctx context.Context, id int) (models.Post, error)
	SelectAllPosts(ctx context.Context) ([]models.Post, error)
	SelectPostsByAuthor(ctx context.Context, authorName string) ([]models.Post, error)
}

type Post struct {
	repo PostRepository
}

func NewPost(repo PostRepository) *Post {
	return &Post{
		repo: repo,
	}
}

func (p *Post) CreatePost(ctx context.Context, text string, authorId int) (int, error) {
	id, err := p.repo.CreatePost(ctx, text, authorId)
	if err != nil {
		return 0, fmt.Errorf("репозиторий постов - вставка записи: %w", err)
	}

	return id, nil
}

func (p *Post) UpdatePost(ctx context.Context, id int, text string) error {
	if err := p.repo.UpdatePost(ctx, id, text); err != nil {
		return fmt.Errorf("репозиторий постов - обновление: %w", err)
	}

	return nil
}

func (p *Post) DeletePost(ctx context.Context, id int) error {
	if err := p.repo.DeletePost(ctx, id); err != nil {
		return fmt.Errorf("репозиторий постов - удаление: %w", err)
	}

	return nil
}

func (p *Post) SelectPostById(ctx context.Context, id int) (models.Post, error) {
	post, err := p.repo.SelectPostById(ctx, id)
	if err != nil {
		return models.Post{}, fmt.Errorf("репозиторий постов - получение записи по id: %w", err)
	}

	return post, nil
}

func (p *Post) SelectAllPosts(ctx context.Context) ([]models.Post, error) {
	posts, err := p.repo.SelectAllPosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("репозиторий постов - получение всех записей: %w", err)
	}

	return posts, nil
}

func (p *Post) SelectPostsByAuthor(ctx context.Context, authorName string) ([]models.Post, error) {
	posts, err := p.repo.SelectPostsByAuthor(ctx, authorName)
	if err != nil {
		return nil, fmt.Errorf("репозиторий постов - получение постов по автору: %w", err)
	}

	return posts, nil
}