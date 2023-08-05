package pasien

import (
	"database/sql"
	"tugas-akhir/internal/repositories/model"
)

type PasienRepo struct {
	db *sql.DB
}

func NewPasienRepositories(db *sql.DB) *PasienRepo {
	return &PasienRepo{db: db}
}

func (p *PasienRepo) FetchDataDiriByID(user_id int64) (model.Pasien, error) {
	var pasien model.Pasien

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
