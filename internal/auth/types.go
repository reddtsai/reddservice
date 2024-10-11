package auth

import (
	"time"
)

type ServiceInput[T any] struct {
	Data T
}

type UserTable struct {
	ID        int64     `json:"id"`
	Account   string    `json:"account"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUser struct {
	Account string `json:"account"`
	Email   string `json:"email"`
}
