package pasien

import (
	"errors"
	"tugas-akhir/internal/repositories/model"
)

func (p *PasienRepo) FetchDokterByPoli(poli_nama string) ([]model.Dokter, error) {
	var dokter []model.Dokter = make([]model.Dokter, 0)

	poliID, _ := p.FetchPoliID(poli_nama)

	var sqlStmt string = `SELECT id, nama from DOKTER where poli_id = ?`

	rows, err := p.db.Query(sqlStmt, poliID)
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

func (p *PasienRepo) FetchDokterIDByJadwalID(jadwal_id int64) (int64, error) {
	var id int64

	var sqlStmt string = `SELECT dokter_id FROM jadwal_dokter WHERE id = ?`

	row := p.db.QueryRow(sqlStmt, jadwal_id)
	err := row.Scan(&id)

	return id, err
}
