package pasien

import (
	"errors"
	"tugas-akhir/internal/repositories/model"
)

func (p *PasienRepo) FetchJadwalDokterByDokterID(dokter_id int64) ([]model.Jadwal, error) {
	var jadwal []model.Jadwal = make([]model.Jadwal, 0)

	var sqlStmt string = `SELECT j.id, j.jadwal_hari, j.jadwal_waktu
		FROM dokter d
		JOIN jadwal_dokter j ON d.id = j.dokter_id
		WHERE d.id = ?`

	rows, err := p.db.Query(sqlStmt, dokter_id)
	if err != nil {
		return nil, errors.New("gagal menampilkan dokter")
	}

	defer rows.Close()

	var dataJadwal model.Jadwal
	for rows.Next() {
		err := rows.Scan(
			&dataJadwal.ID,
			&dataJadwal.JadwalHari,
			&dataJadwal.JadwalWaktu,
		)

		if err != nil {
			return nil, err
		}

		jadwal = append(jadwal, dataJadwal)
	}

	return jadwal, nil
}

func (p *PasienRepo) FetchJadwalByID(id int64) (model.Jadwal, error) {
	var jadwal model.Jadwal

	var sqlStmt string = `SELECT dokter_id, jadwal_hari, jadwal_waktu FROM jadwal_dokter WHERE id = ?`

	row := p.db.QueryRow(sqlStmt, id)
	err := row.Scan(
		&jadwal.Dokter_ID,
		&jadwal.JadwalHari,
		&jadwal.JadwalWaktu,
	)

	return jadwal, err
}
