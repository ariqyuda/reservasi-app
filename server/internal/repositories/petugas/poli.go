package petugas

import "time"

func (prs *PetugasRepo) InsertPoli(nama string) error {

	sqlStmt := `INSERT INTO POLI (nama, created_at) VALUES (?, ?)`

	_, err := prs.db.Exec(sqlStmt, nama, time.Now())
	if err != nil {
		return err
	}

	return err
}
