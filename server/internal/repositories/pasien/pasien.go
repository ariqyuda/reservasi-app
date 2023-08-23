package pasien

import (
	"database/sql"
	"time"
	"tugas-akhir/internal/repositories/model"
	"tugas-akhir/internal/repositories/user"
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
	userRepo := user.NewUserRepositories(p.db)

	id_pasien, _ := userRepo.FetchPasienID(user_id)

	var sqlStmt string = `UPDATE pasien SET nik_pasien = ?, nama = ?, jk_pasien = ?, tgl_lahir_pasien = ?, tmpt_lahir_pasien = ?
	, alamat_pasien = ?, no_hp_pasien = ?, ktp_pasien = ?, updated_at = ? WHERE id = ?`

	_, err := p.db.Exec(sqlStmt, nik, nama, gender, tgl_lahir, tmpt_lahir, alamat, no_hp, ktp_pasien, time.Now(), id_pasien)
	if err != nil {
		return err
	}

	return err
}
