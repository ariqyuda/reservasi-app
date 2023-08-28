package petugas

import (
	"errors"
	"tugas-akhir/internal/repositories/model"
)

func (prs *PetugasRepo) KirimDataLaporan(tanggal_awal, tanggal_akhir string) ([]model.Reservasi, error) {
	var reservasi []model.Reservasi = make([]model.Reservasi, 0)

	var status string = "Selesai"

	var sqlStmt string = `SELECT ps.nik_pasien, ps.nama AS nama_pasien, p.nama AS poli, d.nama AS nama_dokter,
	r.jadwal_tanggal, r.jadwal_hari, r.jadwal_waktu, r.tipe, r.status 
	FROM reservasi r
	JOIN dokter d ON r.dokter_id = d.id
	JOIN poli p ON r.poli_id = p.id
	JOIN pasien ps ON r.user_id = ps.user_id
	WHERE status = ?
	AND 
	jadwal_tanggal BETWEEN ? AND ?`

	rows, err := prs.db.Query(sqlStmt, status, tanggal_awal, tanggal_akhir)
	if err != nil {
		return nil, errors.New("gagal menampilkan data reservasi")
	}

	defer rows.Close()

	var dataReservasi model.Reservasi
	for rows.Next() {
		err := rows.Scan(
			&dataReservasi.NIK,
			&dataReservasi.PasienName,
			&dataReservasi.PoliName,
			&dataReservasi.DokterName,
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
