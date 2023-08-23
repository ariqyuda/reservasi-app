package petugas

import (
	"errors"
	"strings"
	"time"
	"tugas-akhir/internal/repositories/model"
)

func (prs *PetugasRepo) FetchPoli(nama string) (model.Poli, error) {
	var poli model.Poli

	var sqlStmt string = `SELECT nama FROM poli WHERE nama LIKE ?`

	row := prs.db.QueryRow(sqlStmt, nama)
	err := row.Scan(&poli.Name)

	return poli, err
}

func (prs *PetugasRepo) InsertPoli(nama string) error {

	poli, _ := prs.FetchPoli(nama)
	if poli.Name != "" {
		return errors.New("poli sudah terdaftar")
	}

	words := strings.ToLower(nama)
	word := strings.Replace(words, ",", "", -1)
	slug := strings.Replace(word, " ", "-", -1)

	sqlStmt := `INSERT INTO POLI (nama, slug, created_at) VALUES (?, ?, ?)`

	_, err := prs.db.Exec(sqlStmt, nama, slug, time.Now())
	if err != nil {
		return err
	}

	return err
}
