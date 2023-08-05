package petugas

import (
	"errors"
	"time"
	"tugas-akhir/internal/repositories/model"
)

func (prs *PetugasRepo) LihatReservasi() ([]model.Reservasi, error) {
	// var sqlStmt string
	var reservasi []model.Reservasi = make([]model.Reservasi, 0)

	var sqlStmt string = `SELECT d.nama as nama_dokter, p.nama AS poli, ps.nama AS nama_pasien, r.jadwal_tanggal, r.jadwal_hari, r.jadwal_waktu, r.tipe, r.status 
		FROM reservasi r
		JOIN dokter d ON r.dokter_id = d.id
		JOIN poli p ON r.poli_id = p.id
		JOIN pasien ps ON r.user_id = ps.user_id`

	rows, err := prs.db.Query(sqlStmt)
	if err != nil {
		return nil, errors.New("gagal menampilkan data reservasi")
	}

	defer rows.Close()

	var dataReservasi model.Reservasi
	for rows.Next() {
		err := rows.Scan(
			&dataReservasi.DokterName,
			&dataReservasi.PoliName,
			&dataReservasi.PasienName,
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

func (prs *PetugasRepo) VerifikasiReservasi(reservasi_id int64, status string) error {
	var sqlStmt string = `UPDATE reservasi SET status = ?, updated_at = ? WHERE id = ?`

	_, err := prs.db.Exec(sqlStmt, status, time.Now(), reservasi_id)

	if err != nil {
		return err
	}

	return err
}
