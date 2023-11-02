package repositories

import "database/sql"

type PetugasRepo struct {
	db *sql.DB
}

func NewPetugasRepositories(db *sql.DB) *PetugasRepo {
	return &PetugasRepo{db: db}
}

func (prs *PetugasRepo) InsertPetugas(email, nama, password string) error {

	role := "petugas"

	userRepo := NewUserRepositories(prs.db)

	err := userRepo.InsertUser(email, nama, role, password)
	if err != nil {
		return err
	}

	return err
}
