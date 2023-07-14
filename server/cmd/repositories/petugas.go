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

func (prs *PetugasRepo) LihatReservasi() ([]Reservasi, error) {
	// var sqlStmt string
	var reservasi []Reservasi = make([]Reservasi, 0)

	var sqlStmt string = `SELECT d.nama, p.nama, r.jadwal_tanggal, r.jadwal_hari, r.jadwal_waktu, r.tipe, r.status 
		FROM reservasi r
		JOIN dokter d ON r.dokter_id = d.id
		JOIN poli p ON r.poli_id = p.id`

	rows, err := prs.db.Query(sqlStmt)
	if err != nil {
		return nil, errors.New("gagal menampilkan data reservasi")
	}

	defer rows.Close()

	var dataReservasi Reservasi
	for rows.Next() {
		err := rows.Scan(
			&dataReservasi.DokterName,
			&dataReservasi.PoliName,
			&dataReservasi.Tanggal,
			&dataReservasi.Hari,
			&dataReservasi.Waktu,
			&dataReservasi.Tipe,
			&dataReservasi.Status,
		)

		if err != nil {
			return nil, err
		}

		reservasi = append(reservasi, dataReservasi)
	}

	return reservasi, nil
}
