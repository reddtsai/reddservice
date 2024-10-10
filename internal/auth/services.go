package auth

import (
	"database/sql"
	"time"
)

const (
	USER_DEFAULT_PASSWORD = "123456"
)

type service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *service {
	return &service{
		db: db,
	}
}

func (s *service) CreateUser(input ServiceInput[CreateUser]) (int64, error) {
	var userID int64
	stmt, err := s.db.Prepare("INSERT INTO users (NAME, EMAIL, PASSWORD, CREATED_AT, UPDATED_AT) VALUES ($1, $2, $3, $4, $5) RETURNING ID")
	if err != nil {
		return 0, err
	}
	// TODO: hash password
	pwd := USER_DEFAULT_PASSWORD
	now := time.Now().UTC()
	err = stmt.QueryRow(input.Data.Name, input.Data.Email, pwd, now, now).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
