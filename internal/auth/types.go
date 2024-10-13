package auth

import (
	"time"
)

type UserTable struct {
	ID        int64     `json:"id"`
	Account   string    `json:"account"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserInput struct {
	Account string `json:"account"`
	Email   string `json:"email"`
}
