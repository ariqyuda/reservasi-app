package repositories

import (
	"database/sql"
	"errors"
	"tugas-akhir/internal/model"
)

type JadwalRepo struct {
	db *sql.DB
}

func NewJadwalRepositories(db *sql.DB) *JadwalRepo {
	return &JadwalRepo{db: db}
}

func (j *JadwalRepo) InsertJadwal(id_dokter int64, jadwal_hari, jadwal_mulai, jadwal_berakhir string) error {

	jadwal_waktu := jadwal_mulai + " - " + jadwal_berakhir

	sqlStmt := `INSERT INTO jadwal_dokter (dokter_id, jadwal_hari, jadwal_waktu, created_at) VALUES (?, ?, ?, ?)`

	timeRepo := NewTimeRepositories(j.db)
	waktuLokal, _ := timeRepo.SetLocalTime()

	_, err := j.db.Exec(sqlStmt, id_dokter, jadwal_hari, jadwal_waktu, waktuLokal)
	if err != nil {
		return err
	}

	return err
}

func (j *JadwalRepo) FetchJadwalDokter() ([]model.Jadwal, error) {
	var jadwal []model.Jadwal = make([]model.Jadwal, 0)

	var sqlStmt string = `SELECT j.id, d.id, d.nama, j.jadwal_hari, j.jadwal_waktu
		FROM dokter d
		JOIN jadwal_dokter j ON d.id = j.dokter_id`

	rows, err := j.db.Query(sqlStmt)
	if err != nil {
		return nil, errors.New("gagal menampilkan dokter")
	}

	defer rows.Close()

	var dataJadwal model.Jadwal
	for rows.Next() {
		err := rows.Scan(
			&dataJadwal.ID,
			&dataJadwal.Dokter_ID,
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

func (j *JadwalRepo) FetchJadwalDokterByID(dokter_id int64) (model.Jadwal, error) {
	var jadwal model.Jadwal

	var sqlStmt string = `SELECT j.id, d.id, d.nama, j.jadwal_hari, j.jadwal_waktu
		FROM dokter d
		JOIN jadwal_dokter j ON d.id = j.dokter_id
		WHERE d.id = ?`

	row := j.db.QueryRow(sqlStmt, dokter_id)

	err := row.Scan(
		&jadwal.ID,
		&jadwal.Dokter_ID,
		&jadwal.NamaDokter,
		&jadwal.JadwalHari,
		&jadwal.JadwalWaktu,
	)

	if err != nil {
		return jadwal, err
	}

	return jadwal, nil
}

func (j *JadwalRepo) UbahJadwalDokter(reservasi_id int64, jadwal_hari, jadwal_mulai, jadwal_berakhir string) error {

	jadwal_waktu := jadwal_mulai + " - " + jadwal_berakhir

	var sqlStmt string = `UPDATE jadwal_dokter SET jadwal_hari = ?, jadwal_waktu = ?, updated_at = ? WHERE id = ?`

	timeRepo := NewTimeRepositories(j.db)
	waktuLokal, _ := timeRepo.SetLocalTime()

	_, err := j.db.Exec(sqlStmt, jadwal_hari, jadwal_waktu, waktuLokal, reservasi_id)

	if err != nil {
		return err
	}

	return err
}

func (j *JadwalRepo) FetchJadwalDokterByDokterID(dokter_id int64) ([]model.Jadwal, error) {
	var jadwal []model.Jadwal = make([]model.Jadwal, 0)

	var sqlStmt string = `SELECT j.id, j.jadwal_hari, j.jadwal_waktu
		FROM dokter d
		JOIN jadwal_dokter j ON d.id = j.dokter_id
		WHERE d.id = ?`

	rows, err := j.db.Query(sqlStmt, dokter_id)
	if err != nil {
		return nil, errors.New("gagal menampilkan dokter")
	}

	defer rows.Close()

	var dataJadwal model.Jadwal
	for rows.Next() {
		err := rows.Scan(
			&dataJadwal.ID,
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

func (j *JadwalRepo) FetchJadwalByID(id int64) (model.Jadwal, error) {
	var jadwal model.Jadwal

	var sqlStmt string = `SELECT dokter_id, jadwal_hari, jadwal_waktu FROM jadwal_dokter WHERE id = ?`

	row := j.db.QueryRow(sqlStmt, id)
	err := row.Scan(
		&jadwal.Dokter_ID,
		&jadwal.JadwalHari,
		&jadwal.JadwalWaktu,
	)

	return jadwal, err
}

func (j *JadwalRepo) LihatJadwal(user_id int64) ([]model.Reservasi, error) {
	var reservasi []model.Reservasi = make([]model.Reservasi, 0)

	var sqlStmt string = `SELECT p.nama, r.jadwal_tanggal, r.jadwal_hari, r.jadwal_waktu, r.keluhan
	FROM reservasi r
	JOIN dokter d ON r.dokter_id = d.id
	JOIN pasien p ON r.user_id = p.user_id
	WHERE d.user_id = ? AND r.status = 'disetujui'`

	rows, err := j.db.Query(sqlStmt, user_id)
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
			&dataReservasi.Keluhan,
		)

		if err != nil {
			return nil, err
		}

		reservasi = append(reservasi, dataReservasi)
	}

	return reservasi, nil
}
