package repositories

import (
	"database/sql"
	"time"
)

type AdminRepo struct {
	db *sql.DB
}

func NewAdminRepositories(db *sql.DB) *AdminRepo {
	return &AdminRepo{db: db}
}

func (a *AdminRepo) InsertPoli(nama string) error {

	sqlStmt := `INSERT INTO POLI (nama, created_at) VALUES (?, ?)`

	_, err := a.db.Exec(sqlStmt, nama, time.Now())
	if err != nil {
		return err
	}

	return err
}

func (a *AdminRepo) InsertJadwal(id_dokter int64, jadwal_hari, jadwal_waktu string) error {

	sqlStmt := `INSERT INTO jadwal_dokter (dokter_id, jadwal_hari, jadwal_waktu, created_at) VALUES (?, ?, ?, ?)`
	_, err := a.db.Exec(sqlStmt, id_dokter, jadwal_hari, jadwal_waktu, time.Now())
	if err != nil {
		return err
	}

	return err
}

func (a *AdminRepo) InsertDokter(email, nama, password, poli_nama string) error {
	role := "dokter"

	userRepo := NewUserRepositories(a.db)

	err := userRepo.InsertUser(email, nama, role, password)
	if err != nil {
		return err
	}

	idUser, _ := userRepo.FetchUserID(email)
	poliID, _ := userRepo.FetchPoliID(poli_nama)

	sqlStmt := `INSERT INTO dokter (user_id, poli_id, nama, created_at) VALUES (?, ?, ?, ?)`

	_, err = a.db.Exec(sqlStmt, idUser, poliID, nama, time.Now())
	if err != nil {
		return err
	}

	return err
}

func (a *AdminRepo) InsertPetugas(email, nama, password string) error {

	role := "petugas"

	userRepo := NewUserRepositories(a.db)

	err := userRepo.InsertUser(email, nama, role, password)
	if err != nil {
		return err
	}

	return err
}

func (a *AdminRepo) InsertAdmin(email, nama, password string) error {
	role := "admin"

	userRepo := NewUserRepositories(a.db)

	err := userRepo.InsertUser(email, nama, role, password)
	if err != nil {
		return err
	}

	return err
}
