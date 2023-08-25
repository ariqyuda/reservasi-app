package petugas

import (
	"tugas-akhir/internal/repositories/user"
)

func (prs *PetugasRepo) InsertJadwal(id_dokter int64, jadwal_hari, jadwal_mulai, jadwal_berakhir string) error {

	jadwal_waktu := jadwal_mulai + " - " + jadwal_berakhir

	sqlStmt := `INSERT INTO jadwal_dokter (dokter_id, jadwal_hari, jadwal_waktu, created_at) VALUES (?, ?, ?, ?)`

	userRepo := user.NewUserRepositories(prs.db)
	waktuLokal, _ := userRepo.SetLocalTime()

	_, err := prs.db.Exec(sqlStmt, id_dokter, jadwal_hari, jadwal_waktu, waktuLokal)
	if err != nil {
		return err
	}

	return err
}

func (prs *PetugasRepo) UbahJadwalDokter(reservasi_id int64, jadwal_hari, jadwal_mulai, jadwal_berakhir string) error {

	jadwal_waktu := jadwal_mulai + " - " + jadwal_berakhir

	var sqlStmt string = `UPDATE jadwal_dokter SET jadwal_hari = ?, jadwal_waktu = ?, updated_at = ? WHERE id = ?`

	userRepo := user.NewUserRepositories(prs.db)
	waktuLokal, _ := userRepo.SetLocalTime()

	_, err := prs.db.Exec(sqlStmt, jadwal_hari, jadwal_waktu, waktuLokal, reservasi_id)

	if err != nil {
		return err
	}

	return err
}
