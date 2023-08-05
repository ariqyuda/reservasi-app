package auth

import "database/sql"

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepositories(db *sql.DB) *AuthRepo {
	return &AuthRepo{db: db}
}
