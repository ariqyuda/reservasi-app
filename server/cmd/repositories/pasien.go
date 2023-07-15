package repositories

import (
	"database/sql"
	"errors"
	"time"
)

type PasienRepo struct {
	db *sql.DB
}

func NewPasienRepositories(db *sql.DB) *PasienRepo {
	return &PasienRepo{db: db}
}

func (p *PasienRepo) FetchJadwalDokterByPoli(poli_nama string) ([]Jadwal, error) {
	var jadwal []Jadwal = make([]Jadwal, 0)

	userRepo := NewUserRepositories(p.db)

	poliID, _ := userRepo.FetchPoliID(poli_nama)

	var sqlStmt string = `SELECT d.nama, j.jadwal_hari, j.jadwal_waktu
		FROM dokter d
		JOIN jadwal_dokter j ON d.id = j.dokter_id
		WHERE d.poli_id = ?`

	rows, err := p.db.Query(sqlStmt, poliID)
	if err != nil {
		return nil, errors.New("gagal menampilkan dokter")
	}

	defer rows.Close()

	var dataJadwal Jadwal
	for rows.Next() {
		err := rows.Scan(
			&dataJadwal.NamaDokter,
			&dataJadwal.JadwalHari,
			&dataJadwal.JadwalWaktu,
		)

		if err != nil {
			return nil, err
		}

		jadwal = append(jadwal, dataJadwal)
	}

	return jadwal, nil
}

func (p *PasienRepo) FetchJadwalByID(id int64) (Jadwal, error) {
	var jadwal Jadwal

	var sqlStmt string = `SELECT dokter_id, jadwal_hari, jadwal_waktu FROM jadwal_dokter WHERE id = ?`

	row := p.db.QueryRow(sqlStmt, id)
	err := row.Scan(
		&jadwal.Dokter_ID,
		&jadwal.JadwalHari,
		&jadwal.JadwalWaktu,
	)

	return jadwal, err
}

func (p *PasienRepo) FetchDataDiriByID(user_id int64) (Pasien, error) {
	var pasien Pasien

	var sqlStmt string = `SELECT id, user_id, nik_pasien, nama, jk_pasien, tgl_lahir_pasien, tmpt_lahir_pasien, alamat_pasien, no_hp_pasien FROM pasien WHERE user_id = ?`

	row := p.db.QueryRow(sqlStmt, user_id)
	err := row.Scan(
		&pasien.ID,
		&pasien.User_ID,
		&pasien.NIK,
		&pasien.Nama,
		&pasien.Gender,
		&pasien.BornDate,
		&pasien.BornPlace,
		&pasien.Adress,
		&pasien.PhoneNumber,
	)

	return pasien, err
}

func (p *PasienRepo) FetchDokterIDByJadwalID(jadwal_id int64) (int64, error) {
	var id int64

	var sqlStmt string = `SELECT dokter_id FROM jadwal_dokter WHERE id = ?`

	row := p.db.QueryRow(sqlStmt, jadwal_id)
	err := row.Scan(&id)

	return id, err
}

func (p *PasienRepo) FetchPoliIDByDokterID(dokter_id int64) (*int64, error) {
	var id int64

	var sqlStmt string = `SELECT poli_id FROM dokter WHERE id = ?`

	row := p.db.QueryRow(sqlStmt, dokter_id)
	err := row.Scan(&id)

	return &id, err
}

func (p *PasienRepo) ReservasiPribadi(user_id, jadwal_id int64, jadwal_tanggal string) error {
	var jadwal Jadwal
	var pasien Pasien

	pasien, _ = p.FetchDataDiriByID(user_id)
	jadwal, _ = p.FetchJadwalByID(jadwal_id)
	dokterID, _ := p.FetchDokterIDByJadwalID(jadwal_id)
	poliID, _ := p.FetchPoliIDByDokterID(dokterID)
	tipe := "umum"
	status := "menunggu persetujuan"

	var sqlStmt string = `INSERT INTO reservasi (user_id, dokter_id, poli_id, nik_pasien, nama, jk_pasien,
		tgl_lahir_pasien, tmpt_lahir_pasien, alamat_pasien, no_hp_pasien, jadwal_tanggal, jadwal_hari, jadwal_waktu, tipe, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := p.db.Exec(sqlStmt, user_id, dokterID, poliID, pasien.NIK, pasien.Nama, pasien.Gender,
		pasien.BornDate, pasien.BornPlace, pasien.Adress, pasien.PhoneNumber, jadwal_tanggal,
		jadwal.JadwalHari, jadwal.JadwalWaktu, tipe, status, time.Now())

	if err != nil {
		return err
	}

	return err
}

func (p *PasienRepo) FetchReservasiByUserID(user_id int64) ([]Reservasi, error) {
	var reservasi []Reservasi = make([]Reservasi, 0)

	// sqlStmt = `SELECT r.nik_pasien, r.nama, r.jk_pasien, r.tgl_lahir_pasien, r.tmpt_lahir_pasien,
	// 	r.alamat_pasien, r.no_hp_pasien, d.nama, p.nama_poli, r.jadwal_hari, r.jadwal_waktu, r.tipe, r.status
	// 	FROM reservasi r
	// 	JOIN dokter d ON r.dokter_id = d.id
	// 	JOIN poli p ON r.poli_id = p.id`

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

func (p *PasienRepo) FetchLatestReservasiByUserID(user_id int64) (Reservasi, error) {
	var sqlStmt string = `SELECT d.nama, p.nama, r.jadwal_tanggal, r.jadwal_hari, r.jadwal_waktu, r.tipe, r.status 
		FROM reservasi r
		JOIN dokter d ON r.dokter_id = d.id
		JOIN poli p ON r.poli_id = p.id
		WHERE r.id = (SELECT max(r.id) FROM reservasi r)`

	row := p.db.QueryRow(sqlStmt)

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
