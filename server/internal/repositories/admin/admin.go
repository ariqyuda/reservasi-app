package admin

import "database/sql"

type AdminRepo struct {
	db *sql.DB
}

func NewAdminRepositories(db *sql.DB) *AdminRepo {
	return &AdminRepo{db: db}
}
