package model

import "time"

type Pasien struct {
	ID          int64     `db:"id"`
	User_ID     int64     `db:"user_id"`
	Email       string    `db:"email"`
	NIK         string    `db:"nik_pasien"`
	Nama        string    `db:"nama"`
	Gender      string    `db:"jk_pasien"`
	BornDate    string    `db:"tgl_lhr_pasien"`
	BornPlace   string    `db:"tmpt_lhr_pasien"`
	Adress      string    `db:"alamat_pasien"`
	PhoneNumber string    `db:"no_hp_pasien"`
	KTP         string    `json:"ktp_pasien"`
	CreatedAt   time.Time `db:"created_at"`
}
