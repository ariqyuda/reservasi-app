package model

type Jadwal struct {
	ID          int64  `db:"id"`
	Dokter_ID   int64  `db:"dokter_id"`
	NamaDokter  string `db:"nama_dokter"`
	JadwalHari  string `json:"jadwal_hari"`
	JadwalWaktu string `json:"jadwal_waktu"`
}
