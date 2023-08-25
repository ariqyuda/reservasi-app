package petugas

import (
	"tugas-akhir/internal/repositories/pasien"
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

	pasienRepo := pasien.NewPasienRepositories(prs.db)
	poliID, _ := pasienRepo.FetchPoliIDByNama(poli_nama)

	sqlStmt := `INSERT INTO dokter (user_id, poli_id, nama, created_at) VALUES (?, ?, ?, ?)`

	waktuLokal, _ := userRepo.SetLocalTime()

	_, err = prs.db.Exec(sqlStmt, idUser, poliID, nama, waktuLokal)
	if err != nil {
		return err
	}

	return err
}
