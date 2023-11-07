package repositories

import (
	"database/sql"
	"errors"
	"tugas-akhir/internal/model"
)

type ReservasiRepo struct {
	db *sql.DB
}

func NewReservasiRepositories(db *sql.DB) *ReservasiRepo {
	return &ReservasiRepo{db: db}
}

func (r *ReservasiRepo) ReservasiPribadi(user_id, jadwal_id int64, jadwal_tanggal, keluhan string) error {
	var jadwal model.Jadwal
	var pasien model.Pasien

	pasienRepo := NewPasienRepositories(r.db)
	jadwalRepo := NewJadwalRepositories(r.db)
	poliRepo := NewPoliRepositories(r.db)
	dokterRepo := NewDokterRepositories(r.db)

	pasien, _ = pasienRepo.FetchDataDiriByID(user_id)
	jadwal, _ = jadwalRepo.FetchJadwalByID(jadwal_id)
	dokterID, _ := dokterRepo.FetchDokterIDByJadwalID(jadwal_id)
	poliID, _ := poliRepo.FetchPoliIDByDokterID(dokterID)
	tipe := "umum"
	status := "menunggu persetujuan"

	var sqlStmt string = `INSERT INTO reservasi (user_id, dokter_id, poli_id, nik_pasien, nama, jk_pasien,
		tgl_lahir_pasien, tmpt_lahir_pasien, alamat_pasien, no_hp_pasien, jadwal_tanggal, jadwal_hari, jadwal_waktu, tipe, status, keluhan, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	timeRepo := NewTimeRepositories(r.db)
	waktuLokal, _ := timeRepo.SetLocalTime()

	_, err := r.db.Exec(sqlStmt, user_id, dokterID, poliID, pasien.NIK, pasien.Nama, pasien.Gender,
		pasien.BornDate, pasien.BornPlace, pasien.Adress, pasien.PhoneNumber, jadwal_tanggal,
		jadwal.JadwalHari, jadwal.JadwalWaktu, tipe, status, keluhan, waktuLokal)

	if err != nil {
		return err
	}

	return err
}

func (p *PasienRepo) FetchReservasiByUserID(user_id int64) ([]model.Reservasi, error) {
	var reservasi []model.Reservasi = make([]model.Reservasi, 0)

	var sqlStmt string = `SELECT d.nama, p.nama, r.nama, r.jadwal_tanggal, r.jadwal_hari, r.jadwal_waktu, r.tipe, r.status 
		FROM reservasi r
		JOIN dokter d ON r.dokter_id = d.id
		JOIN poli p ON r.poli_id = p.id
		WHERE r.user_id = ?
		ORDER BY r.created_at DESC`

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

func (p *PasienRepo) FetchLatestReservasiByUserID(user_id int64) (model.Reservasi, error) {
	var sqlStmt string = `SELECT d.nama, p.nama, r.jadwal_tanggal, r.jadwal_hari, r.jadwal_waktu, r.tipe, r.status 
		FROM reservasi r
		JOIN dokter d ON r.dokter_id = d.id
		JOIN poli p ON r.poli_id = p.id
		WHERE r.user_id = ?
		ORDER BY r.created_at DESC LIMIT 1`

	row := p.db.QueryRow(sqlStmt, user_id)

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

func (r *ReservasiRepo) LihatReservasi() ([]model.Reservasi, error) {
	// var sqlStmt string
	var reservasi []model.Reservasi = make([]model.Reservasi, 0)

	var sqlStmt string = `SELECT r.id, d.nama as nama_dokter, p.nama AS poli, r.nama AS nama_pasien, r.jadwal_tanggal, r.jadwal_hari, r.jadwal_waktu, r.tipe, r.status 
		FROM reservasi r
		JOIN dokter d ON r.dokter_id = d.id
		JOIN poli p ON r.poli_id = p.id
		LEFT JOIN pasien ps ON r.user_id = ps.user_id
		WHERE NOT r.status = 'Selesai'
		ORDER BY r.created_at DESC`

	rows, err := r.db.Query(sqlStmt)
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

func (prs *PetugasRepo) VerifikasiReservasi(reservasi_id int64, status, alasan_verifikasi string) error {
	var sqlStmt string = `UPDATE reservasi SET status = ?, alasan_verifikasi = ?, updated_at = ? WHERE id = ?`

	timeRepo := NewTimeRepositories(prs.db)
	waktuLokal, _ := timeRepo.SetLocalTime()

	_, err := prs.db.Exec(sqlStmt, status, alasan_verifikasi, waktuLokal, reservasi_id)

	if err != nil {
		return err
	}

	return err
}

func (prs *PetugasRepo) ReservasiPasien(user_id, jadwal_id int64, jadwal_tanggal, nik_pasien, nama_pasien, jk_pasien,
	tgl_lahir_pasien, tmpt_lahir_pasien, alamat_pasien, no_hp_pasien, keluhan string) error {

	var jadwal model.Jadwal

	jadwalRepo := NewJadwalRepositories(prs.db)
	dokterRepo := NewDokterRepositories(prs.db)
	poliRepo := NewPoliRepositories(prs.db)

	jadwal, _ = jadwalRepo.FetchJadwalByID(jadwal_id)
	dokterID, _ := dokterRepo.FetchDokterIDByJadwalID(jadwal_id)
	poliID, _ := poliRepo.FetchPoliIDByDokterID(dokterID)
	tipe := "umum"
	status := "menunggu persetujuan"

	var sqlStmt string = `INSERT INTO reservasi (user_id, dokter_id, poli_id, nik_pasien, nama, jk_pasien,
		tgl_lahir_pasien, tmpt_lahir_pasien, alamat_pasien, no_hp_pasien, jadwal_tanggal, jadwal_hari, jadwal_waktu, tipe, status, keluhan, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	timeRepo := NewTimeRepositories(prs.db)
	waktuLokal, _ := timeRepo.SetLocalTime()

	_, err := prs.db.Exec(sqlStmt, user_id, dokterID, poliID, nik_pasien, nama_pasien, jk_pasien,
		tgl_lahir_pasien, tmpt_lahir_pasien, alamat_pasien, no_hp_pasien, jadwal_tanggal,
		jadwal.JadwalHari, jadwal.JadwalWaktu, tipe, status, keluhan, waktuLokal)

	if err != nil {
		return err
	}

	return err
}
