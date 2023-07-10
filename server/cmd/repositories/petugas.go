package repositories

import "database/sql"

type PetugasRepo struct {
	db *sql.DB
}

func NewPetugasRepositories(db *sql.DB) *PetugasRepo {
	return &PetugasRepo{db: db}
}

func (p *PetugasRepo) LihatReservasi() ([]Reservasi, error) {
	// var sqlStmt string
	var reservasi []Reservasi = make([]Reservasi, 0)

	return reservasi, nil
}
