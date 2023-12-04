package repositories

import (
	"database/sql"
	"errors"
)

type PetugasRepo struct {
	db *sql.DB
}

func NewPetugasRepositories(db *sql.DB) *PetugasRepo {
	return &PetugasRepo{db: db}
}

func (prs *PetugasRepo) InsertPetugas(email, nama, password string) error {

	// check input
	if email == "" || nama == "" || password == "" {
		return errors.New("input tidak boleh kosong")
	}

	role := "petugas"
	status := "aktif"

	userRepo := NewUserRepositories(prs.db)

	err := userRepo.InsertUser(email, nama, role, password, status)
	if err != nil {
		return err
	}

	return err
}
