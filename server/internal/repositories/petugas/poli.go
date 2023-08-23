package petugas

import (
	"strings"
	"time"
)

func (prs *PetugasRepo) InsertPoli(nama string) error {

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
