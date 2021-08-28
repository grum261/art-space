package models

import (
	"art_space/internal"
)

type User struct {
	Id       int
	Username string
	Password string
	Dates    Dates
	Bio      string
	Avatar   string
}

func (u *User) Validate() error {
	if u.Username == "" {
		return internal.NewErrorf("никнейм автора поста не может быть пустым", internal.ErrorCodeInvalidArgument)
	}

	if u.Password == "" {
		return internal.NewErrorf("пароль не должен быть пустым", internal.ErrorCodeInvalidArgument)
	}

	return nil
}
