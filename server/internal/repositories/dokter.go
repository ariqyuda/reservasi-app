package repositories

import (
	"database/sql"
	"errors"
	"tugas-akhir/internal/model"
)

type DokterRepo struct {
	db *sql.DB
}

func NewDokterRepositories(db *sql.DB) *DokterRepo {
	return &DokterRepo{db: db}
}

func (d *DokterRepo) FetchDokter() ([]model.Dokter, error) {
	var dokter []model.Dokter = make([]model.Dokter, 0)

	var sqlStmt = `SELECT id, nama FROM dokter`

	rows, err := d.db.Query(sqlStmt)
	if err != nil {
		return nil, errors.New("gagal menampilkan dokter")
	}

	defer rows.Close()

	var dataDokter model.Dokter
	for rows.Next() {
		err := rows.Scan(
			&dataDokter.ID,
			&dataDokter.Nama,
		)

		if err != nil {
			return nil, err
		}

		dokter = append(dokter, dataDokter)
	}

	return dokter, nil
}

func (d *DokterRepo) InsertDokter(email, nama, password, str_dokter, sip_dokter, poli_nama, status_dokter string) error {
	role := "dokter"
	status := "aktif"

	userRepo := NewUserRepositories(d.db)

	err := userRepo.InsertUser(email, nama, role, password, status)
	if err != nil {
		return err
	}

	idUser, _ := userRepo.FetchUserID(email)

	poliRepo := NewPoliRepositories(d.db)
	poliID, _ := poliRepo.FetchPoliIDByNama(poli_nama)

	sqlStmt := `INSERT INTO dokter (user_id, poli_id, nama, str_dokter, sip_dokter, status, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`

	timeRepo := NewTimeRepositories(d.db)
	waktuLokal, _ := timeRepo.SetLocalTime()

	_, err = d.db.Exec(sqlStmt, idUser, poliID, nama, str_dokter, sip_dokter, status_dokter, waktuLokal)
	if err != nil {
		return err
	}

	return err
}

func (d *DokterRepo) FetchDokterByPoliNama(slug string) ([]model.Dokter, error) {
	var dokter []model.Dokter = make([]model.Dokter, 0)

	poliRepo := NewPoliRepositories(d.db)

	poliID, _ := poliRepo.FetchPoliIDBySlug(slug)

	var sqlStmt string = `SELECT d.id, d.nama, p.nama 
	from dokter d
	JOIN poli P ON d.poli_id = p.id 
	where poli_id = ?`

	rows, err := d.db.Query(sqlStmt, poliID)
	if err != nil {
		return nil, errors.New("gagal menampilkan dokter")
	}

	defer rows.Close()

	var dataDokter model.Dokter
	for rows.Next() {
		err := rows.Scan(
			&dataDokter.ID,
			&dataDokter.Nama,
			&dataDokter.PoliName,
		)

		if err != nil {
			return nil, err
		}

		dokter = append(dokter, dataDokter)
	}

	return dokter, nil
}

func (d *DokterRepo) FetchDokterByID(dokter_id int64) (model.Dokter, error) {
	var dokter model.Dokter

	var sqlStmt string = `SELECT id, nama, poli_id FROM dokter WHERE id = ?`

	row := d.db.QueryRow(sqlStmt, dokter_id)
	err := row.Scan(
		&dokter.ID,
		&dokter.Nama,
		&dokter.PoliID,
	)

	return dokter, err
}

func (d *DokterRepo) FetchDokterIDByJadwalID(jadwal_id int64) (int64, error) {
	var id int64

	var sqlStmt string = `SELECT dokter_id FROM jadwal_dokter WHERE id = ?`

	row := d.db.QueryRow(sqlStmt, jadwal_id)
	err := row.Scan(&id)

	return id, err
}

func (d *DokterRepo) UbahStatusDokter(dokter_id int64, status string) error {
	var sqlStmt string = `UPDATE dokter SET status = ?, updated_at = ? WHERE id = ?`

	timeRepo := NewTimeRepositories(d.db)
	waktuLokal, _ := timeRepo.SetLocalTime()

	_, err := d.db.Exec(sqlStmt, status, waktuLokal, dokter_id)

	if err != nil {
		return err
	}

	return err
}
