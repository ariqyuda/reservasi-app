package repositories

import (
	"database/sql"
	"errors"
)

type Pasien struct {
	ID          int64  `db:"id"`
	User_ID     int64  `db:"user_id"`
	Email       string `db:"email"`
	NIK         string `db:"nik_pasien"`
	Nama        string `db:"nama"`
	Gender      string `db:"jk_pasien"`
	BornDate    string `db:"tgl_lhr_pasien"`
	BornPlace   string `db:"tmpt_lhr_pasien"`
	Adress      string `db:"alamat_pasien"`
	PhoneNumber string `db:"no_hp_pasien"`
	KTP         string `json:"ktp_pasien"`
}

type PasienRepo struct {
	db *sql.DB
}

func NewPasienRepositories(db *sql.DB) *PasienRepo {
	return &PasienRepo{db: db}
}

func (p *PasienRepo) FetchDataDiriByID(user_id int64) (Pasien, error) {
	var pasien Pasien

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

	id_pasien, _ := p.FetchPasienID(user_id)

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

func (p *PasienRepo) FetchDataPasien() ([]Pasien, error) {
	var user []Pasien = make([]Pasien, 0)

	var sqlStmt string = `SELECT u.id, u.email, u.nama, p.nik_pasien
	FROM users u
	JOIN pasien p ON u.id = p.user_id`

	rows, err := p.db.Query(sqlStmt)
	if err != nil {
		return nil, errors.New("gagal menampilkan data user")
	}

	defer rows.Close()

	var dataUser Pasien
	for rows.Next() {
		err := rows.Scan(
			&dataUser.ID,
			&dataUser.Email,
			&dataUser.Nama,
			&dataUser.NIK,
		)

		if err != nil {
			return nil, err
		}

		user = append(user, dataUser)
	}

	return user, nil
}

func (p *PasienRepo) FetchPasienByNIK(nik string) (Pasien, error) {
	var pasien Pasien

	var sqlStmt string = `SELECT nik_pasien FROM pasien WHERE nik_pasien = ?`

	row := p.db.QueryRow(sqlStmt, nik)
	err := row.Scan(&pasien.NIK)

	return pasien, err
}

func (p *PasienRepo) FetchPasienID(user_id int64) (*int64, error) {
	var sqlStmt string = `SELECT id FROM pasien WHERE user_id = ?`
	var id int64

	row := p.db.QueryRow(sqlStmt, user_id)
	err := row.Scan(&id)

	return &id, err
}
