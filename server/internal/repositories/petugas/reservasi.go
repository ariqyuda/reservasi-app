package petugas

import (
	"errors"
	"tugas-akhir/internal/repositories/model"
	"tugas-akhir/internal/repositories/pasien"
	"tugas-akhir/internal/repositories/user"
)

func (prs *PetugasRepo) LihatReservasi() ([]model.Reservasi, error) {
	// var sqlStmt string
	var reservasi []model.Reservasi = make([]model.Reservasi, 0)

	var sqlStmt string = `SELECT r.id, d.nama as nama_dokter, p.nama AS poli, r.nama AS nama_pasien, r.jadwal_tanggal, r.jadwal_hari, r.jadwal_waktu, r.tipe, r.status 
		FROM reservasi r
		JOIN dokter d ON r.dokter_id = d.id
		JOIN poli p ON r.poli_id = p.id
		LEFT JOIN pasien ps ON r.user_id = ps.user_id
		WHERE NOT r.status = 'Selesai'
		ORDER BY r.created_at DESC`

	rows, err := prs.db.Query(sqlStmt)
	if err != nil {
		return nil, errors.New("gagal menampilkan data reservasi")
	}

	defer rows.Close()

	var dataReservasi model.Reservasi
	for rows.Next() {
		err := rows.Scan(
			&dataReservasi.ID,
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

	userRepo := user.NewUserRepositories(prs.db)
	waktuLokal, _ := userRepo.SetLocalTime()

	_, err := prs.db.Exec(sqlStmt, status, waktuLokal, reservasi_id)

	if err != nil {
		return err
	}

	return err
}

func (prs *PetugasRepo) ReservasiPasien(user_id, jadwal_id int64, jadwal_tanggal, nik_pasien, nama_pasien, jk_pasien,
	tgl_lahir_pasien, tmpt_lahir_pasien, alamat_pasien, no_hp_pasien string) error {

	var jadwal model.Jadwal

	pasienRepo := pasien.NewPasienRepositories(prs.db)

	jadwal, _ = pasienRepo.FetchJadwalByID(jadwal_id)
	dokterID, _ := pasienRepo.FetchDokterIDByJadwalID(jadwal_id)
	poliID, _ := pasienRepo.FetchPoliIDByDokterID(dokterID)
	tipe := "umum"
	status := "menunggu persetujuan"

	var sqlStmt string = `INSERT INTO reservasi (user_id, dokter_id, poli_id, nik_pasien, nama, jk_pasien,
		tgl_lahir_pasien, tmpt_lahir_pasien, alamat_pasien, no_hp_pasien, jadwal_tanggal, jadwal_hari, jadwal_waktu, tipe, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	userRepo := user.NewUserRepositories(prs.db)
	waktuLokal, _ := userRepo.SetLocalTime()

	_, err := prs.db.Exec(sqlStmt, user_id, dokterID, poliID, nik_pasien, nama_pasien, jk_pasien,
		tgl_lahir_pasien, tmpt_lahir_pasien, alamat_pasien, no_hp_pasien, jadwal_tanggal,
		jadwal.JadwalHari, jadwal.JadwalWaktu, tipe, status, waktuLokal)

	if err != nil {
		return err
	}

	return err
}
