package pasien

import (
	"errors"
	"tugas-akhir/internal/repositories/model"
)

func (p *PasienRepo) FetchDokterByPoliNama(slug string) ([]model.Dokter, error) {
	var dokter []model.Dokter = make([]model.Dokter, 0)

	poliID, _ := p.FetchPoliIDBySlug(slug)

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

func (p *PasienRepo) FetchDokterByID(dokter_id int64) (model.Dokter, error) {
	var dokter model.Dokter

	var sqlStmt string = `SELECT id, nama, poli_id FROM dokter WHERE id = ?`

	row := p.db.QueryRow(sqlStmt, dokter_id)
	err := row.Scan(
		&dokter.ID,
		&dokter.Nama,
		&dokter.PoliID,
	)

	return dokter, err
}

func (p *PasienRepo) FetchDokterIDByJadwalID(jadwal_id int64) (int64, error) {
	var id int64

	var sqlStmt string = `SELECT dokter_id FROM jadwal_dokter WHERE id = ?`

	row := p.db.QueryRow(sqlStmt, jadwal_id)
	err := row.Scan(&id)

	return id, err
}
