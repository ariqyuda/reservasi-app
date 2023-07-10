package repositories

import (
	"database/sql"
	"errors"
)

type DokterRepo struct {
	db *sql.DB
}

func NewDokterRepositories(db *sql.DB) *DokterRepo {
	return &DokterRepo{db: db}
}

func (d *DokterRepo) LihatJadwal(user_id int64) ([]Reservasi, error) {
	var reservasi []Reservasi = make([]Reservasi, 0)

	var sqlStmt string = `SELECT`
	_, err := d.db.Query(sqlStmt)
	if err != nil {
		return nil, errors.New("gagal menampilkan jadwal")
	}

	return reservasi, nil
}
