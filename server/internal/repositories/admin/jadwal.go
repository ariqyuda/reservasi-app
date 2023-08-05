package admin

import "time"

func (a *AdminRepo) InsertJadwal(id_dokter int64, jadwal_hari, jadwal_waktu string) error {

	sqlStmt := `INSERT INTO jadwal_dokter (dokter_id, jadwal_hari, jadwal_waktu, created_at) VALUES (?, ?, ?, ?)`
	_, err := a.db.Exec(sqlStmt, id_dokter, jadwal_hari, jadwal_waktu, time.Now())
	if err != nil {
		return err
	}

	return err
}
