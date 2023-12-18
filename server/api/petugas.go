package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Poli struct {
	Name string `json:"nama_poli"`
}

type Dokter struct {
	DokterID     int64  `json:"id_dokter"`
	Email        string `json:"email"`
	Name         string `json:"nama"`
	Password     string `json:"password"`
	Poli         string `json:"nama_poli"`
	STRDokter    string `json:"str_dokter"`
	SIPDokter    string `json:"sip_dokter"`
	StatusDokter string `json:"status"`
}

type JadwalDokter struct {
	DokterID      int64  `json:"id_dokter"`
	Hari          string `json:"jadwal_hari"`
	WaktuMulai    string `json:"jadwal_mulai"`
	WaktuBerakhir string `json:"jadwal_berakhir"`
}

type UbahPoli struct {
	ID   int64  `db:"id"`
	Nama string `db:"nama"`
}

type UbahJadwalDokter struct {
	JadwalID      int64  `json:"id"`
	Hari          string `json:"jadwal_hari"`
	WaktuMulai    string `json:"jadwal_mulai"`
	WaktuBerakhir string `json:"jadwal_berakhir"`
}

type VerifReservasi struct {
	ID               int64  `json:"id"`
	Status           string `json:"status"`
	AlasanVerifikasi string `json:"alasan_verifikasi"`
}

type StatusDokter struct {
	DokterID int64  `json:"id_dokter"`
	Status   string `json:"status"`
}

type DataLaporan struct {
	WaktuAwal  string `json:"waktu_awal"`
	WaktuAkhir string `json:"waktu_akhir"`
}

type ReservasiPasien struct {
	JadwalID       int64  `json:"id_jadwal"`
	Jadwal_Tanggal string `json:"jadwal_tanggal"`
	Nama           string `json:"nama_pasien"`
	NIK            string `json:"nik_pasien"`
	Gender         string `json:"jk_pasien"`
	BornDate       string `json:"tgl_lahir_pasien"`
	BornPlace      string `json:"tmpt_lahir_pasien"`
	Adress         string `json:"alamat_pasien"`
	PhoneNumber    string `json:"no_hp_pasien"`
	Keluhan        string `json:"keluhan"`
}

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

type FetchDataSuccessResponse struct {
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination"`
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
	err = api.poliRepo.InsertPoli(poli.Name)
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
	err = api.dokterRepo.InsertDokter(dokter.Email, dokter.Name, dokter.Password, dokter.STRDokter, dokter.SIPDokter, dokter.StatusDokter, dokter.Poli)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	insertDokter := Dokter{
		Email:     dokter.Email,
		Name:      dokter.Name,
		Poli:      dokter.Poli,
		STRDokter: dokter.STRDokter,
		SIPDokter: dokter.SIPDokter,
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
	err = api.jadwalRepo.InsertJadwal(jadwal.DokterID, jadwal.Hari, jadwal.WaktuMulai, jadwal.WaktuBerakhir)
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

func (api *API) fetchDokter(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	dokter, err := api.dokterRepo.FetchDokter()
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(ReservasiErrorResponse{Error: err.Error()})
		}
	}()
	fetchDokterResponse := FetchDataSuccessResponse{
		Message: "success",
		Data:    dokter,
	}

	json.NewEncoder(w).Encode(fetchDokterResponse)
}

func (api *API) fetchJadwalDokter(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	reservasi, err := api.jadwalRepo.FetchJadwalDokter()
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(ReservasiErrorResponse{Error: err.Error()})
		}
	}()
	fetchJadwalResponse := FetchDataSuccessResponse{
		Message: "success",
		Data:    reservasi,
	}

	json.NewEncoder(w).Encode(fetchJadwalResponse)
}

func (api *API) lihatReservasiUser(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	page, err := strconv.Atoi(req.URL.Query().Get("page"))

	reservasi, pagination, err := api.reservasiRepo.LihatReservasi(page)
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(ReservasiErrorResponse{Error: err.Error()})
		}
	}()
	lihatReservasiResponse := FetchDataSuccessResponse{
		Message:    "success",
		Data:       reservasi,
		Pagination: pagination,
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
	err = api.reservasiRepo.VerifikasiReservasi(bodyRequest.ID, bodyRequest.Status, bodyRequest.AlasanVerifikasi)
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

func (api *API) kirimDataLaporan(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	var bodyRequest DataLaporan
	err := json.NewDecoder(req.Body).Decode(&bodyRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	reservasi, err := api.reservasiRepo.DataLaporanReservasi(bodyRequest.WaktuAwal, bodyRequest.WaktuAkhir)

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

func (api *API) ubahPoli(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	var bodyRequest UbahPoli
	err := json.NewDecoder(req.Body).Decode(&bodyRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.poliRepo.UbahPoli(bodyRequest.ID, bodyRequest.Nama)
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
	err = api.jadwalRepo.UbahJadwalDokter(bodyRequest.JadwalID, bodyRequest.Hari, bodyRequest.WaktuMulai, bodyRequest.WaktuBerakhir)
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

func (api *API) reservasiPasien(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	userID := req.Context().Value("id").(int64)

	var reservasi ReservasiPasien
	err := json.NewDecoder(req.Body).Decode(&reservasi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.reservasiRepo.ReservasiPasien(int64(userID), reservasi.JadwalID, reservasi.Jadwal_Tanggal, reservasi.NIK, reservasi.Nama, reservasi.Gender,
		reservasi.BornDate, reservasi.BornPlace, reservasi.Adress, reservasi.PhoneNumber, reservasi.Keluhan)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	reservasiResponse := PasienSuccessResponse{
		Message: "reservasi berhasil",
		Data:    reservasi,
	}

	json.NewEncoder(w).Encode(reservasiResponse)
}

func (api *API) ubahStatusDokter(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	var bodyRequest StatusDokter
	err := json.NewDecoder(req.Body).Decode(&bodyRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.dokterRepo.UbahStatusDokter(bodyRequest.DokterID, bodyRequest.Status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	statusDokterResponse := VerifikasiResponse{
		Message: "Berhasil merubah status dokter",
	}

	json.NewEncoder(w).Encode(statusDokterResponse)
}

func (api *API) ubahDataDokter(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	var bodyRequest Dokter
	err := json.NewDecoder(req.Body).Decode(&bodyRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.dokterRepo.UbahDataDokter(bodyRequest.DokterID, bodyRequest.Name, bodyRequest.STRDokter, bodyRequest.SIPDokter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	statusDokterResponse := VerifikasiResponse{
		Message: "Berhasil memperbarui data dokter",
	}

	json.NewEncoder(w).Encode(statusDokterResponse)
}
