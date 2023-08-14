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
	ID            int64  `json:"id"`
	Hari          string `json:"jadwal_hari"`
	WaktuMulai    string `json:"jadwal_mulai"`
	WaktuBerakhir string `json:"jadwal_berakhir"`
}

type FetchUserSuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type InsertSuccessResponse struct {
	Message string   `json:"message"`
	Data    Register `json:"data"`
}

func (api *API) lihatDataUser(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	userRole := req.URL.Query().Get("role")

	reservasi, err := api.usersRepo.FetchDataUserByRole(userRole)
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(ReservasiErrorResponse{Error: err.Error()})
		}
	}()
	fetchUserResponse := FetchUserSuccessResponse{
		Message: "success",
		Data:    reservasi,
	}

	json.NewEncoder(w).Encode(fetchUserResponse)
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
