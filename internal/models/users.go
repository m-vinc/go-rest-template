package models

import (
	"errors"
	"mpj/pkg/ent"
	"time"
)

type User struct {
	*ent.User
}

type UserCreator struct {
	Username    string     `json:"username" validate:"required"`
	Firstname   *string    `json:"first_name" validate:"omitempty,required"`
	Lastname    *string    `json:"last_name" validate:"omitempty,required"`
	DateOfBirth *time.Time `json:"date_of_birth" validate:"omitempty"`
	Description *string    `json:"description" validate:"omitempty"`
}

type UserQuery struct {
	Username *string `json:"user_hash" validate:"omitempty"`
}

var (
	ErrUsernameAlreadyTaken = errors.New("user: this user already exist in the database, unlucky ?")
	ErrUserNoGetters        = errors.New("user: cannot get a user without at least one getter attribute")
	ErrUserNotFound         = errors.New("user: not found")
)
