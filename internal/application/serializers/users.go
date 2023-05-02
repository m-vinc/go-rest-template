package serializers

import (
	"mpj/internal/models"
	"mpj/pkg/openapi/v1"
)

func User(u *models.User) *openapi.User {
	m := &openapi.User{
		Username:    &u.Username,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		DateOfBirth: u.DateOfBirth,
		Description: u.Description,
	}

	return m
}

func Users(us []*models.User) *[]openapi.User {
	ms := make([]openapi.User, len(us))

	for i := 0; i < len(us); i++ {
		ms[i] = *User(us[i])
	}

	return &ms
}
