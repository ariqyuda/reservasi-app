package api

import (
	"encoding/json"
	"net/http"
)

type Reservasi struct {
	ID         int64  `db:"id"`
	DokterName string `db:"dokter_name"`
	PoliName   string `db:"poli_name"`
	Hari       string `json:"jadwal_hari"`
	Waktu      string `json:"jadwal_waktu"`
	Type       string `db:"tipe"`
	Status     string `db:"status"`
}

type Jadwal struct {
	ID           int64  `db:"jadwal_id"`
	Jadwal_Hari  string `json:"jadwal_hari"`
	Jadwal_Waktu string `json:"jadwal_waktu"`
}

type ReservasiSuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ReservasiErrorResponse struct {
	Error string `json:"error"`
}

type PasienErrorResponse struct {
	Error string `json:"error"`
}

func (api *API) reservasiPribadi(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	userID := req.Context().Value("id").(int64)

	var jadwal Jadwal
	err := json.NewDecoder(req.Body).Decode(&jadwal)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.pasienRepo.ReservasiPribadi(int64(userID), jadwal.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	reservasiResponse := ReservasiSuccessResponse{
		Message: "reservasi berhasil",
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
	reservasiResponse := ReservasiSuccessResponse{
		Message: "success",
		Data:    dataReservasi,
	}

	encoder.Encode(reservasiResponse)
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

	reservasiResponse := ReservasiSuccessResponse{
		Message: "success",
		Data:    hasilReservasi,
	}

	encoder.Encode(reservasiResponse)
}
