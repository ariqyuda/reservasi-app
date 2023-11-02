package repositories

import (
	"database/sql"
	"tugas-akhir/internal/model"
)

type PasienRepo struct {
	db *sql.DB
}

func NewPasienRepositories(db *sql.DB) *PasienRepo {
	return &PasienRepo{db: db}
}

func (p *PasienRepo) FetchDataDiriByID(user_id int64) (model.Pasien, error) {
	var pasien model.Pasien

	var sqlStmt string = `SELECT nik_pasien, nama, jk_pasien, tgl_lahir_pasien, tmpt_lahir_pasien, alamat_pasien, no_hp_pasien, ktp_pasien FROM pasien WHERE user_id = ?`

	row := p.db.QueryRow(sqlStmt, user_id)
	err := row.Scan(
		&pasien.NIK,
		&pasien.Nama,
		&pasien.Gender,
		&pasien.BornDate,
		&pasien.BornPlace,
		&pasien.Adress,
		&pasien.PhoneNumber,
		&pasien.KTP,
	)

	return pasien, err
}

func (p *PasienRepo) UbahDataDiri(user_id int64, nik, nama, gender, tgl_lahir, tmpt_lahir, alamat, no_hp, ktp_pasien string) error {

	userRepo := NewUserRepositories(p.db)
	id_pasien, _ := userRepo.FetchPasienID(user_id)

	var sqlStmt string = `UPDATE pasien SET nik_pasien = ?, nama = ?, jk_pasien = ?, tgl_lahir_pasien = ?, tmpt_lahir_pasien = ?
	, alamat_pasien = ?, no_hp_pasien = ?, ktp_pasien = ?, updated_at = ? WHERE id = ?`

	timeRepo := NewTimeRepositories(p.db)
	waktuLokal, _ := timeRepo.SetLocalTime()

	_, err := p.db.Exec(sqlStmt, nik, nama, gender, tgl_lahir, tmpt_lahir, alamat, no_hp, ktp_pasien, waktuLokal, id_pasien)
	if err != nil {
		return err
	}

	return err
}
