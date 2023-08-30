package pasien

import (
	"errors"
	"tugas-akhir/internal/repositories/model"
)

func (p *PasienRepo) FetchPoliIDByNama(nama string) (*int64, error) {
	var sqlStmt string
	var id int64

	sqlStmt = `SELECT id FROM poli WHERE nama = ?`

	row := p.db.QueryRow(sqlStmt, nama)
	err := row.Scan(&id)

	return &id, err
}

func (p *PasienRepo) FetchPoliIDBySlug(slug string) (*int64, error) {
	var sqlStmt string
	var id int64

	sqlStmt = `SELECT id FROM poli WHERE slug = ?`

	row := p.db.QueryRow(sqlStmt, slug)
	err := row.Scan(&id)

	return &id, err
}

func (p *PasienRepo) FetchPoli() ([]model.Poli, error) {
	var poli []model.Poli = make([]model.Poli, 0)

	var sqlStmt string = `SELECT id, nama, slug FROM poli`

	rows, err := p.db.Query(sqlStmt)
	if err != nil {
		return nil, errors.New("gagal menampilkan poli")
	}

	defer rows.Close()

	var dataPoli model.Poli
	for rows.Next() {
		err := rows.Scan(
			&dataPoli.ID,
			&dataPoli.Name,
			&dataPoli.Slug,
		)

		if err != nil {
			return nil, err
		}

		poli = append(poli, dataPoli)
	}

	return poli, nil
}

func (p *PasienRepo) FetchPoliIDByDokterID(dokter_id int64) (*int64, error) {
	var id int64

	var sqlStmt string = `SELECT poli_id FROM dokter WHERE id = ?`

	row := p.db.QueryRow(sqlStmt, dokter_id)
	err := row.Scan(&id)

	return &id, err
}
