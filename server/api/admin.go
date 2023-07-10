package api

import (
	"encoding/json"
	"net/http"
)

type AdminErrorResponse struct {
	Error string `json:"error"`
}

type AdminResponse struct {
	Error string `json:"error"`
}

type Dokter struct {
	Email    string `json:"email"`
	Name     string `json:"nama"`
	Password string `json:"password"`
	Poli     string `json:"nama_poli"`
}

type Register struct {
	Email    string `json:"email"`
	Name     string `json:"nama"`
	Password string `json:"password"`
}

type Poli struct {
	Name string `json:"nama_poli"`
}

type JadwalDokter struct {
	ID    int64  `json:"id"`
	Hari  string `json:"jadwal_hari"`
	Waktu string `json:"jadwal_waktu"`
}

type InsertDokterSuccessResponse struct {
	Message string `json:"message"`
	Data    Dokter `json:"data"`
}

type InsertSuccessResponse struct {
	Message string   `json:"message"`
	Data    Register `json:"data"`
}

type InsertJadwalDokterSuccessResponse struct {
	Message string       `json:"message"`
	Data    JadwalDokter `json:"data"`
}

type InsertPoliSuccessResponse struct {
	Message string `json:"message"`
	Data    Poli   `json:"data"`
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
	err = api.adminRepo.InsertDokter(dokter.Email, dokter.Name, dokter.Password, dokter.Poli)
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

func (api *API) insertPetugas(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	var user Register
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.adminRepo.InsertPetugas(user.Email, user.Name, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	insertUser := Register{
		Email: user.Email,
		Name:  user.Name,
	}

	insertDokterResponse := InsertSuccessResponse{
		Message: "insert success",
		Data:    insertUser,
	}

	json.NewEncoder(w).Encode(insertDokterResponse)
}

func (api *API) insertAdmin(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	var user Register
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.adminRepo.InsertAdmin(user.Email, user.Name, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	insertUser := Register{
		Email: user.Email,
		Name:  user.Name,
	}

	insertDokterResponse := InsertSuccessResponse{
		Message: "insert success",
		Data:    insertUser,
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
	err = api.adminRepo.InsertJadwal(jadwal.ID, jadwal.Hari, jadwal.Waktu)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	insertJadwalDokter := JadwalDokter{
		ID:    jadwal.ID,
		Hari:  jadwal.Hari,
		Waktu: jadwal.Waktu,
	}

	insertDokterResponse := InsertJadwalDokterSuccessResponse{
		Message: "insert success",
		Data:    insertJadwalDokter,
	}

	json.NewEncoder(w).Encode(insertDokterResponse)
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
	err = api.adminRepo.InsertPoli(poli.Name)
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
