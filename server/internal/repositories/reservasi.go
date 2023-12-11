package repositories

import (
	"database/sql"
	"errors"
)

type Reservasi struct {
	ID               int64  `db:"id"`
	NIKPasien        string `db:"nik_pasien"`
	NamaPasien       string `db:"nama"`
	Gender           string `db:"jk_pasien"`
	BornDate         string `db:"tgl_lhr_pasien"`
	BornPlace        string `db:"tmpt_lhr_pasien"`
	Adress           string `db:"alamat_pasien"`
	PhoneNumber      string `db:"no_hp_pasien"`
	DokterName       string `db:"dokter_name"`
	PoliName         string `db:"nama_poli"`
	Tanggal          string `db:"jadwal_tanggal"`
	Hari             string `json:"jadwal_hari"`
	Waktu            string `json:"jadwal_waktu"`
	Tipe             string `db:"tipe"`
	Status           string `db:"status"`
	Keluhan          string `db:"keluhan"`
	AlasanVerifikasi string `db:"alasan_verifikasi"`
}

type ReservasiRepo struct {
	db *sql.DB
}

func NewReservasiRepositories(db *sql.DB) *ReservasiRepo {
	return &ReservasiRepo{db: db}
}

func (r *ReservasiRepo) ReservasiPribadi(user_id, jadwal_id int64, jadwal_tanggal, keluhan string) error {
	var jadwal Jadwal
	var pasien Pasien

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

func (p *PasienRepo) FetchReservasiByUserID(user_id int64) ([]Reservasi, error) {
	var reservasi []Reservasi = make([]Reservasi, 0)

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

	var dataReservasi Reservasi
	for rows.Next() {
		err := rows.Scan(
			&dataReservasi.DokterName,
			&dataReservasi.PoliName,
			&dataReservasi.NamaPasien,
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

func (p *PasienRepo) FetchLatestReservasiByUserID(user_id int64) (Reservasi, error) {
	var sqlStmt string = `SELECT d.nama, p.nama, r.jadwal_tanggal, r.jadwal_hari, r.jadwal_waktu, r.tipe, r.status 
		FROM reservasi r
		JOIN dokter d ON r.dokter_id = d.id
		JOIN poli p ON r.poli_id = p.id
		WHERE r.user_id = ?
		ORDER BY r.created_at DESC LIMIT 1`

	row := p.db.QueryRow(sqlStmt, user_id)

	var reservasi Reservasi
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

func (r *ReservasiRepo) LihatReservasi() ([]Reservasi, error) {
	// var sqlStmt string
	var reservasi []Reservasi = make([]Reservasi, 0)

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

	var dataReservasi Reservasi
	for rows.Next() {
		err := rows.Scan(
			&dataReservasi.ID,
			&dataReservasi.DokterName,
			&dataReservasi.PoliName,
			&dataReservasi.NamaPasien,
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

func (r *ReservasiRepo) VerifikasiReservasi(reservasi_id int64, status, alasan_verifikasi string) error {
	var sqlStmt string = `UPDATE reservasi SET status = ?, alasan_verifikasi = ?, updated_at = ? WHERE id = ?`

	timeRepo := NewTimeRepositories(r.db)
	waktuLokal, _ := timeRepo.SetLocalTime()

	_, err := r.db.Exec(sqlStmt, status, alasan_verifikasi, waktuLokal, reservasi_id)

	if err != nil {
		return err
	}

	return err
}

func (r *ReservasiRepo) ReservasiPasien(user_id, jadwal_id int64, jadwal_tanggal, nik_pasien, nama_pasien, jk_pasien,
	tgl_lahir_pasien, tmpt_lahir_pasien, alamat_pasien, no_hp_pasien, keluhan string) error {

	// check input
	if nik_pasien == "" || nama_pasien == "" || jk_pasien == "" || tgl_lahir_pasien == "" || tmpt_lahir_pasien == "" ||
		alamat_pasien == "" || no_hp_pasien == "" || keluhan == "" {
		return errors.New("data tidak boleh kosong")
	}

	var jadwal Jadwal

	jadwalRepo := NewJadwalRepositories(r.db)
	dokterRepo := NewDokterRepositories(r.db)
	poliRepo := NewPoliRepositories(r.db)

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

	_, err := r.db.Exec(sqlStmt, user_id, dokterID, poliID, nik_pasien, nama_pasien, jk_pasien,
		tgl_lahir_pasien, tmpt_lahir_pasien, alamat_pasien, no_hp_pasien, jadwal_tanggal,
		jadwal.JadwalHari, jadwal.JadwalWaktu, tipe, status, keluhan, waktuLokal)

	if err != nil {
		return err
	}

	return err
}

func (r *ReservasiRepo) DataLaporanReservasi(tanggal_awal, tanggal_akhir string) ([]Reservasi, error) {
	var reservasi []Reservasi = make([]Reservasi, 0)

	var status string = "Selesai"

	var sqlStmt string = `SELECT ps.nik_pasien, ps.nama AS nama_pasien, p.nama AS poli, d.nama AS nama_dokter,
	r.jadwal_tanggal, r.jadwal_hari, r.jadwal_waktu, r.tipe, r.status, r.alasan_verifikasi 
	FROM reservasi r
	JOIN dokter d ON r.dokter_id = d.id
	JOIN poli p ON r.poli_id = p.id
	JOIN pasien ps ON r.user_id = ps.user_id
	WHERE r.status = ?
	AND 
	r.jadwal_tanggal BETWEEN ? AND ?`

	rows, err := r.db.Query(sqlStmt, status, tanggal_awal, tanggal_akhir)
	if err != nil {
		return nil, errors.New("gagal menampilkan data reservasi")
	}

	defer rows.Close()

	var dataReservasi Reservasi
	for rows.Next() {
		err := rows.Scan(
			&dataReservasi.NIKPasien,
			&dataReservasi.NamaPasien,
			&dataReservasi.PoliName,
			&dataReservasi.DokterName,
			&dataReservasi.Tanggal,
			&dataReservasi.Hari,
			&dataReservasi.Waktu,
			&dataReservasi.Tipe,
			&dataReservasi.Status,
			&dataReservasi.AlasanVerifikasi,
		)

		if err != nil {
			return nil, err
		}

		reservasi = append(reservasi, dataReservasi)
	}

	return reservasi, nil
}

func (r *ReservasiRepo) LihatJadwalReservasi(user_id int64) ([]Reservasi, error) {
	var reservasi []Reservasi = make([]Reservasi, 0)

	var sqlStmt string = `SELECT p.nama, r.jadwal_tanggal, r.jadwal_hari, r.jadwal_waktu, r.keluhan
	FROM reservasi r
	JOIN dokter d ON r.dokter_id = d.id
	JOIN pasien p ON r.user_id = p.user_id
	WHERE d.user_id = ? AND r.status = 'disetujui'`

	rows, err := r.db.Query(sqlStmt, user_id)
	if err != nil {
		return nil, errors.New("gagal menampilkan jadwal")
	}

	defer rows.Close()

	var dataReservasi Reservasi
	for rows.Next() {
		err := rows.Scan(
			&dataReservasi.NamaPasien,
			&dataReservasi.Tanggal,
			&dataReservasi.Hari,
			&dataReservasi.Waktu,
			&dataReservasi.Keluhan,
		)

		if err != nil {
			return nil, err
		}

		reservasi = append(reservasi, dataReservasi)
	}

	return reservasi, nil
}
