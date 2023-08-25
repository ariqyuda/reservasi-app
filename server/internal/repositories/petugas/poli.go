package petugas

import (
	"errors"
	"strings"
	"tugas-akhir/internal/repositories/model"
	"tugas-akhir/internal/repositories/user"
)

func (prs *PetugasRepo) FetchPoliByNama(nama string) (model.Poli, error) {
	var poli model.Poli

	var sqlStmt string = `SELECT nama FROM poli WHERE nama LIKE ?`

	row := prs.db.QueryRow(sqlStmt, nama)
	err := row.Scan(&poli.Name)

	return poli, err
}

func (prs *PetugasRepo) SetSlugByName(nama string) (string, error) {
	// membuat kata menjadi huruf kecil
	words := strings.ToLower(nama)

	// menghilangkan tanda koma
	words = strings.Replace(words, ",", "", -1)

	// menggantikan spasi dengan tanda strip
	slug := strings.Replace(words, " ", "-", -1)

	return slug, nil
}

func (prs *PetugasRepo) InsertPoli(nama string) error {

	poli, _ := prs.FetchPoliByNama(nama)
	if poli.Name != "" {
		return errors.New("poli sudah terdaftar")
	}

	slug, _ := prs.SetSlugByName(nama)

	sqlStmt := `INSERT INTO POLI (nama, slug, created_at) VALUES (?, ?, ?)`

	userRepo := user.NewUserRepositories(prs.db)
	waktuLokal, _ := userRepo.SetLocalTime()

	_, err := prs.db.Exec(sqlStmt, nama, slug, waktuLokal)
	if err != nil {
		return err
	}

	return err
}

func (prs *PetugasRepo) UbahPoli(poli_id int64, nama string) error {

	slug, _ := prs.SetSlugByName(nama)

	var sqlStmt string = `UPDATE poli SET nama = ?, slug = ?, updated_at = ? WHERE id = ?`

	userRepo := user.NewUserRepositories(prs.db)
	waktuLokal, _ := userRepo.SetLocalTime()

	_, err := prs.db.Exec(sqlStmt, nama, slug, waktuLokal, poli_id)

	if err != nil {
		return err
	}

	return err
}
