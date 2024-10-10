package auth

import (
	"time"
)

type ServiceInput[T any] struct {
	Data T
}

type UserTable struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
