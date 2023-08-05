package admin

import "time"

func (a *AdminRepo) InsertPoli(nama string) error {

	sqlStmt := `INSERT INTO POLI (nama, created_at) VALUES (?, ?)`

	_, err := a.db.Exec(sqlStmt, nama, time.Now())
	if err != nil {
		return err
	}

	return err
}
