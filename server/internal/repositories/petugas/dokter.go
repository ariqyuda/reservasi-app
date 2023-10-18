package petugas

import (
	"errors"
	"tugas-akhir/internal/repositories/model"
	"tugas-akhir/internal/repositories/pasien"
	"tugas-akhir/internal/repositories/user"
)

func (prs *PetugasRepo) FetchDokter() ([]model.Dokter, error) {
	var dokter []model.Dokter = make([]model.Dokter, 0)

	var sqlStmt = `SELECT id, nama FROM dokter`

	rows, err := prs.db.Query(sqlStmt)
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
