package model

import "time"

type Dokter struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Email     string    `db:"email"`
	Nama      string    `db:"nama_dokter"`
	STRDokter string    `db:"str_dokter"`
	SIPDokter string    `db:"sip_dokter"`
	PoliID    int64     `db:"poli_id"`
	PoliName  string    `db:"nama_poli"`
	CreatedAt time.Time `db:"created_at"`
}
