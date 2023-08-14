package api

import (
	"encoding/json"
	"net/http"
)

type InsertPoliSuccessResponse struct {
	Message string `json:"message"`
	Data    Poli   `json:"data"`
}

type InsertDokterSuccessResponse struct {
	Message string `json:"message"`
	Data    Dokter `json:"data"`
}

type InsertJadwalDokterSuccessResponse struct {
	Message string       `json:"message"`
	Data    JadwalDokter `json:"data"`
}

type ReservasiResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type VerifReservasi struct {
	ID     int64  `db:"id"`
	Status string `db:"status"`
}

type VerifikasiResponse struct {
	Message string `json:"message"`
}

func (api *API) insertPoli(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	var poli Poli
	err := json.NewDecoder(req.Body).Decode(&poli)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.petugasRepo.InsertPoli(poli.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	insertPoli := Poli{
		Name: poli.Name,
	}

	insertDokterResponse := InsertPoliSuccessResponse{
		Message: "insert success",
		Data:    insertPoli,
	}

	json.NewEncoder(w).Encode(insertDokterResponse)
}

func (api *API) insertDokter(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	var dokter Dokter
	err := json.NewDecoder(req.Body).Decode(&dokter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.petugasRepo.InsertDokter(dokter.Email, dokter.Name, dokter.Password, dokter.Poli)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	insertDokter := Dokter{
		Email: dokter.Email,
		Name:  dokter.Name,
		Poli:  dokter.Poli,
	}

	insertDokterResponse := InsertDokterSuccessResponse{
		Message: "insert success",
		Data:    insertDokter,
	}

	json.NewEncoder(w).Encode(insertDokterResponse)
}

func (api *API) insertJadwalDokter(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	var jadwal JadwalDokter
	err := json.NewDecoder(req.Body).Decode(&jadwal)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.petugasRepo.InsertJadwal(jadwal.ID, jadwal.Hari, jadwal.WaktuMulai, jadwal.WaktuBerakhir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	insertJadwalDokter := JadwalDokter{
		ID:            jadwal.ID,
		Hari:          jadwal.Hari,
		WaktuMulai:    jadwal.WaktuMulai,
		WaktuBerakhir: jadwal.WaktuBerakhir,
	}

	insertJadwalDokterResponse := InsertJadwalDokterSuccessResponse{
		Message: "insert success",
		Data:    insertJadwalDokter,
	}

	json.NewEncoder(w).Encode(insertJadwalDokterResponse)
}

func (api *API) lihatReservasiUser(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	reservasi, err := api.petugasRepo.LihatReservasi()
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(ReservasiErrorResponse{Error: err.Error()})
		}
	}()
	reservasiResponse := ReservasiSuccessResponse{
		Message: "success",
		Data:    reservasi,
	}

	json.NewEncoder(w).Encode(reservasiResponse)
}

func (api *API) verifikasiReservasi(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	var bodyRequest VerifReservasi
	err := json.NewDecoder(req.Body).Decode(&bodyRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.petugasRepo.VerifikasiReservasi(bodyRequest.ID, bodyRequest.Status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	verifikasiResponse := VerifikasiResponse{
		Message: "Berhasil merubah status reservasi",
	}

	json.NewEncoder(w).Encode(verifikasiResponse)
}

func (api *API) ubahJadwalDokter(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	var bodyRequest JadwalDokter
	err := json.NewDecoder(req.Body).Decode(&bodyRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.petugasRepo.UbahJadwalDokter(bodyRequest.ID, bodyRequest.Hari, bodyRequest.WaktuMulai, bodyRequest.WaktuBerakhir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	verifikasiResponse := VerifikasiResponse{
		Message: "Berhasil merubah jadwal dokter",
	}

	json.NewEncoder(w).Encode(verifikasiResponse)
}
