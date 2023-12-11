package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Jadwal struct {
	ID             int64  `json:"id"`
	Jadwal_Tanggal string `json:"jadwal_tanggal"`
	Jadwal_Hari    string `json:"jadwal_hari"`
	Jadwal_Waktu   string `json:"jadwal_waktu"`
}

type PasienSuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type FethcDokterResponse struct {
	Message string      `json:"message"`
	Poli    string      `json:"poli"`
	Data    interface{} `json:"data"`
}

type FetchJadwalSuccessResponse struct {
	Message string      `json:"message"`
	Dokter  interface{} `json:"data_dokter"`
	Data    interface{} `json:"data"`
}

type ReservasiErrorResponse struct {
	Error string `json:"error"`
}

type PasienErrorResponse struct {
	Error string `json:"error"`
}

func (api *API) lihatDataDiri(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	userID := req.Context().Value("id").(int64)

	pasien, err := api.pasienRepo.FetchDataDiriByID(userID)
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(ReservasiErrorResponse{Error: err.Error()})
		}
	}()
	fetchDataDiriResponse := PasienSuccessResponse{
		Message: "success",
		Data:    pasien,
	}

	json.NewEncoder(w).Encode(fetchDataDiriResponse)
}

func (api *API) ubahDataDiri(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	userID := req.Context().Value("id").(int64)

	var bodyRequest Pasien
	err := json.NewDecoder(req.Body).Decode(&bodyRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.pasienRepo.UbahDataDiri(userID, bodyRequest.NIK, bodyRequest.Nama, bodyRequest.Gender, bodyRequest.BornDate, bodyRequest.BornPlace,
		bodyRequest.Adress, bodyRequest.PhoneNumber, bodyRequest.KTP)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	ubahDataDiriResponse := PasienSuccessResponse{
		Message: "Berhasil melakukan perubahan data diri",
		Data:    bodyRequest,
	}

	json.NewEncoder(w).Encode(ubahDataDiriResponse)
}

func (api *API) lihatPoli(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	poli, err := api.poliRepo.FetchPoli()
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(ReservasiErrorResponse{Error: err.Error()})
		}
	}()
	fetchPoliResponse := PasienSuccessResponse{
		Message: "success",
		Data:    poli,
	}

	json.NewEncoder(w).Encode(fetchPoliResponse)
}

func (api *API) lihatDokter(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	poli := req.URL.Query().Get("poli")
	page, err := strconv.Atoi(req.URL.Query().Get("page"))

	reservasi, err := api.dokterRepo.FetchDokterByPoliNama(poli, page)
	poliName, err := api.poliRepo.FetchPoliNameBySlug(poli)
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(ReservasiErrorResponse{Error: err.Error()})
		}
	}()
	fetchJadwalResponse := FethcDokterResponse{
		Message: "success",
		Poli:    poliName,
		Data:    reservasi,
	}

	json.NewEncoder(w).Encode(fetchJadwalResponse)
}

func (api *API) lihatJadwalDokter(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	id_dokter, err := strconv.Atoi(req.URL.Query().Get("id"))
	page, err := strconv.Atoi(req.URL.Query().Get("page"))

	dokter, err := api.dokterRepo.FetchDokterByID(int64(id_dokter))
	reservasi, err := api.jadwalRepo.FetchJadwalDokterByDokterID(int64(id_dokter), page)
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(ReservasiErrorResponse{Error: err.Error()})
		}
	}()
	fetchJadwalResponse := FetchJadwalSuccessResponse{
		Message: "success",
		Dokter:  dokter,
		Data:    reservasi,
	}

	json.NewEncoder(w).Encode(fetchJadwalResponse)
}

func (api *API) reservasiPribadi(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	userID := req.Context().Value("id").(int64)

	var reservasi ReservasiPasien
	err := json.NewDecoder(req.Body).Decode(&reservasi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.reservasiRepo.ReservasiPribadi(int64(userID), reservasi.JadwalID, reservasi.Jadwal_Tanggal, reservasi.Keluhan)
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

func (api *API) lihatReservasi(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	userID := req.Context().Value("id").(int64)

	dataReservasi, err := api.pasienRepo.FetchReservasiByUserID(userID)
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(ReservasiErrorResponse{Error: err.Error()})
		}
	}()
	lihatReservasiResponse := PasienSuccessResponse{
		Message: "success",
		Data:    dataReservasi,
	}

	encoder.Encode(lihatReservasiResponse)
}

func (api *API) lihatHasilReservasi(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	userID := req.Context().Value("id").(int64)

	encoder := json.NewEncoder(w)
	hasilReservasi, err := api.pasienRepo.FetchLatestReservasiByUserID(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	lihatHasilResponse := PasienSuccessResponse{
		Message: "success",
		Data:    hasilReservasi,
	}

	encoder.Encode(lihatHasilResponse)
}
