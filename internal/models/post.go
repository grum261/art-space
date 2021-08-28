package models

import (
	"art_space/internal"
	"time"
)

type Post struct {
	Id       int
	Text     string
	Dates    Dates
	Comments []Comment
	Author   Author
}

type Comment struct {
	Id     int
	Text   string
	Dates  Dates
	Author Author
}

type Author struct {
	Id     int
	Name   string
	Avatar string
}

type Dates struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (d *Dates) Validate() error {
	if !d.CreatedAt.IsZero() && !d.UpdatedAt.IsZero() && d.UpdatedAt.Before(d.CreatedAt) {
		return internal.NewErrorf("дата обновления не может быть меньше даты создания", internal.ErrorCodeInvalidArgument)
	}

	return nil
}

func (p *Post) Validate() error {
	if p.Text == "" {
		return internal.NewErrorf("текст поста не может быть пустым", internal.ErrorCodeInvalidArgument)
	}

	if p.Author.Name == "" {
		return internal.NewErrorf("никнейм автора поста не может быть пустым", internal.ErrorCodeInvalidArgument)
	}

	return nil
}

func (c *Comment) Validate() error {
	if c.Text == "" {
		return internal.NewErrorf("текст комментария не может быть пустым", internal.ErrorCodeInvalidArgument)
	}

	if c.Author.Name == "" {
		return internal.NewErrorf("никнейм автора комментария не может быть пустым", internal.ErrorCodeInvalidArgument)
	}

	return nil
}
