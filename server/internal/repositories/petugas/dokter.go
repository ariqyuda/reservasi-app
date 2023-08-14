package petugas

import (
	"time"
	"tugas-akhir/internal/repositories/user"
)

func (prs *PetugasRepo) InsertDokter(email, nama, password, poli_nama string) error {
	role := "dokter"

	userRepo := user.NewUserRepositories(prs.db)

	err := userRepo.InsertUser(email, nama, role, password)
	if err != nil {
		return err
	}

	idUser, _ := userRepo.FetchUserID(email)
	poliID, _ := userRepo.FetchPoliID(poli_nama)

	sqlStmt := `INSERT INTO dokter (user_id, poli_id, nama, created_at) VALUES (?, ?, ?, ?)`

	_, err = prs.db.Exec(sqlStmt, idUser, poliID, nama, time.Now())
	if err != nil {
		return err
	}

	return err
}
