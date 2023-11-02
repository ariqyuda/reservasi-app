package model

import "time"

type Reservasi struct {
	ID          int64     `db:"id"`
	NIK         string    `db:"nik_pasien"`
	PasienName  string    `db:"pasien_name"`
	Gender      string    `db:"jk_pasien"`
	BornDate    string    `db:"tgl_lhr_pasien"`
	BornPlace   string    `db:"tmpt_lhr_pasien"`
	Adress      string    `db:"alamat_pasien"`
	PhoneNumber string    `db:"no_hp_pasien"`
	DokterName  string    `db:"dokter_name"`
	PoliName    string    `db:"nama_poli"`
	Tanggal     string    `db:"jadwal_tanggal"`
	Hari        string    `json:"jadwal_hari"`
	Waktu       string    `json:"jadwal_waktu"`
	Tipe        string    `db:"tipe"`
	Status      string    `db:"status"`
	Keluhan     string    `db:"keluhan"`
	CreatedAt   time.Time `db:"created_at"`
}
