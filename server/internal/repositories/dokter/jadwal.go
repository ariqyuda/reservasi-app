package dokter

import (
	"database/sql"
	"errors"
	"tugas-akhir/internal/repositories/model"
)

type DokterRepo struct {
	db *sql.DB
}

func NewDokterRepositories(db *sql.DB) *DokterRepo {
	return &DokterRepo{db: db}
}

func (d *DokterRepo) LihatJadwal(user_id int64) ([]model.Reservasi, error) {
	var reservasi []model.Reservasi = make([]model.Reservasi, 0)

	var sqlStmt string = `SELECT p.nama, r.jadwal_tanggal, r.jadwal_hari, r.jadwal_waktu
	FROM reservasi r
	JOIN dokter d ON r.dokter_id = d.id
	JOIN pasien p ON r.user_id = p.user_id
	WHERE d.user_id = ? AND r.status = 'disetujui'`

	rows, err := d.db.Query(sqlStmt, user_id)
	if err != nil {
		return nil, errors.New("gagal menampilkan jadwal")
	}

	defer rows.Close()

	var dataReservasi model.Reservasi
	for rows.Next() {
		err := rows.Scan(
			&dataReservasi.PasienName,
			&dataReservasi.Tanggal,
			&dataReservasi.Hari,
			&dataReservasi.Waktu,
		)

		if err != nil {
			return nil, err
		}

		reservasi = append(reservasi, dataReservasi)
	}

	return reservasi, nil
}
