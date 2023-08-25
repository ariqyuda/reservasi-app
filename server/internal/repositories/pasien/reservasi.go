package pasien

import (
	"errors"
	"tugas-akhir/internal/repositories/model"
	"tugas-akhir/internal/repositories/user"
)

func (p *PasienRepo) ReservasiPribadi(user_id, jadwal_id int64, jadwal_tanggal string) error {
	var jadwal model.Jadwal
	var pasien model.Pasien

	pasien, _ = p.FetchDataDiriByID(user_id)
	jadwal, _ = p.FetchJadwalByID(jadwal_id)
	dokterID, _ := p.FetchDokterIDByJadwalID(jadwal_id)
	poliID, _ := p.FetchPoliIDByDokterID(dokterID)
	tipe := "umum"
	status := "menunggu persetujuan"

	var sqlStmt string = `INSERT INTO reservasi (user_id, dokter_id, poli_id, nik_pasien, nama, jk_pasien,
		tgl_lahir_pasien, tmpt_lahir_pasien, alamat_pasien, no_hp_pasien, jadwal_tanggal, jadwal_hari, jadwal_waktu, tipe, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	userRepo := user.NewUserRepositories(p.db)
	waktuLokal, _ := userRepo.SetLocalTime()

	_, err := p.db.Exec(sqlStmt, user_id, dokterID, poliID, pasien.NIK, pasien.Nama, pasien.Gender,
		pasien.BornDate, pasien.BornPlace, pasien.Adress, pasien.PhoneNumber, jadwal_tanggal,
		jadwal.JadwalHari, jadwal.JadwalWaktu, tipe, status, waktuLokal)

	if err != nil {
		return err
	}

	return err
}

func (p *PasienRepo) FetchReservasiByUserID(user_id int64) ([]model.Reservasi, error) {
	var reservasi []model.Reservasi = make([]model.Reservasi, 0)

	var sqlStmt string = `SELECT d.nama, p.nama, r.jadwal_tanggal, r.jadwal_hari, r.jadwal_waktu, r.tipe, r.status 
		FROM reservasi r
		JOIN dokter d ON r.dokter_id = d.id
		JOIN poli p ON r.poli_id = p.id
		WHERE r.user_id = ?`

	rows, err := p.db.Query(sqlStmt, user_id)
	if err != nil {
		return nil, errors.New("gagal menampilkan data reservasi")
	}

	defer rows.Close()

	var dataReservasi model.Reservasi
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

func (p *PasienRepo) FetchLatestReservasiByUserID(user_id int64) (model.Reservasi, error) {
	var sqlStmt string = `SELECT d.nama, p.nama, r.jadwal_tanggal, r.jadwal_hari, r.jadwal_waktu, r.tipe, r.status 
		FROM reservasi r
		JOIN dokter d ON r.dokter_id = d.id
		JOIN poli p ON r.poli_id = p.id
		WHERE r.id = (SELECT max(r.id) FROM reservasi r)`

	row := p.db.QueryRow(sqlStmt)

	var reservasi model.Reservasi
	err := row.Scan(
		&reservasi.DokterName,
		&reservasi.PoliName,
		&reservasi.Tanggal,
		&reservasi.Hari,
		&reservasi.Waktu,
		&reservasi.Tipe,
		&reservasi.Status,
	)

	if err != nil {
		return reservasi, err
	}

	return reservasi, nil
}
