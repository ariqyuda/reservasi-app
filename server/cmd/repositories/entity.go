package repositories

import "time"

type User struct {
	ID        int64     `db:"id"`
	Email     string    `db:"email"`
	Name      string    `db:"nama"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	Token     string    `db:"token"`
}

type Pasien struct {
	ID          int64     `db:"id"`
	User_ID     int64     `db:"user_id"`
	NIK         string    `db:"nik_pasien"`
	Nama        string    `db:"nama"`
	Gender      string    `db:"jk_pasien"`
	BornDate    string    `db:"tgl_lhr_pasien"`
	BornPlace   string    `db:"tmpt_lhr_pasien"`
	Adress      string    `db:"alamat_pasien"`
	PhoneNumber string    `db:"no_hp_pasien"`
	CreatedAt   time.Time `db:"created_at"`
}

type Dokter struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Role      string    `json:"role"`
	Hari      string    `json:"jadwal_hari"`
	Waktu     string    `json:"jadwal_waktu"`
	PoliID    int64     `db:"poli_id"`
	CreatedAt time.Time `db:"created_at"`
}

type Poli struct {
	ID        int64     `db:"id"`
	Name      string    `db:"nama_poli"`
	CreatedAt time.Time `db:"created_at"`
}

type Jadwal struct {
	ID           int64  `db:"id"`
	dokter_ID    int64  `db:"dokter_id"`
	Jadwal_Hari  string `json:"jadwal_hari"`
	Jadwal_Waktu string `json:"jadwal_waktu"`
}

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
	Hari        string    `json:"jadwal_hari"`
	Waktu       string    `json:"jadwal_waktu"`
	Tipe        string    `db:"tipe"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
}
