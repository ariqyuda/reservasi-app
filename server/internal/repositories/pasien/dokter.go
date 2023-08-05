package pasien

func (p *PasienRepo) FetchDokterIDByJadwalID(jadwal_id int64) (int64, error) {
	var id int64

	var sqlStmt string = `SELECT dokter_id FROM jadwal_dokter WHERE id = ?`

	row := p.db.QueryRow(sqlStmt, jadwal_id)
	err := row.Scan(&id)

	return id, err
}
