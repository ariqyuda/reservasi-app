package petugas

import "tugas-akhir/internal/repositories/model"

func (prs *PetugasRepo) KirimDataLaporan(tanggal_awal, tanggal_akhir string) ([]model.Reservasi, error) {
	var reservasi []model.Reservasi = make([]model.Reservasi, 0)

	return reservasi, nil
}
