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

type UbahJadwalDokter struct {
	ReservasiID   int64  `json:"id"`
	Hari          string `json:"jadwal_hari"`
	WaktuMulai    string `json:"jadwal_mulai"`
	WaktuBerakhir string `json:"jadwal_berakhir"`
}

type ReservasiResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type VerifReservasi struct {
	ID     int64  `db:"id"`
	Status string `db:"status"`
}

type UbahPoli struct {
	ID   int64  `db:"id"`
	Nama string `db:"nama"`
}

type VerifikasiResponse struct {
	Message string `json:"message"`
}

type UbahPoliResponse struct {
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

	insertPoliResponse := InsertPoliSuccessResponse{
		Message: "insert success",
		Data:    insertPoli,
	}

	json.NewEncoder(w).Encode(insertPoliResponse)
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
	err = api.petugasRepo.InsertJadwal(jadwal.DokterID, jadwal.Hari, jadwal.WaktuMulai, jadwal.WaktuBerakhir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	insertJadwalDokter := JadwalDokter{
		DokterID:      jadwal.DokterID,
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
	lihatReservasiResponse := ReservasiResponse{
		Message: "success",
		Data:    reservasi,
	}

	json.NewEncoder(w).Encode(lihatReservasiResponse)
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

func (api *API) ubahPoli(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	var bodyRequest UbahPoli
	err := json.NewDecoder(req.Body).Decode(&bodyRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.petugasRepo.UbahPoli(bodyRequest.ID, bodyRequest.Nama)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	ubahPoliResponse := UbahPoliResponse{
		Message: "Berhasil merubah poli",
	}

	json.NewEncoder(w).Encode(ubahPoliResponse)
}

func (api *API) ubahJadwalDokter(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	var bodyRequest UbahJadwalDokter
	err := json.NewDecoder(req.Body).Decode(&bodyRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.petugasRepo.UbahJadwalDokter(bodyRequest.ReservasiID, bodyRequest.Hari, bodyRequest.WaktuMulai, bodyRequest.WaktuBerakhir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	ubahJadwalResponse := VerifikasiResponse{
		Message: "Berhasil merubah jadwal dokter",
	}

	json.NewEncoder(w).Encode(ubahJadwalResponse)
}
