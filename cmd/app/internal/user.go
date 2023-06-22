package main

import "go-flow2sqlite/cmd/app/internal/models"

type User struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

func (u User) ToModel() models.User {
	return models.User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}

}

type conn interface {
	Connect()
	Do()
}
