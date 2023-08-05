package petugas

import "database/sql"

type PetugasRepo struct {
	db *sql.DB
}

func NewPetugasRepositories(db *sql.DB) *PetugasRepo {
	return &PetugasRepo{db: db}
}
