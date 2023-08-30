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

type Petugas struct {
	Email    string `json:"email"`
	Name     string `json:"nama"`
	Password string `json:"password"`
}

type FetchUserSuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type InsertPetugasSuccessResponse struct {
	Message string  `json:"message"`
	Data    Petugas `json:"data"`
}

func (api *API) lihatDataUser(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	userRole := req.URL.Query().Get("role")

	userData, err := api.usersRepo.FetchDataUserByRole(userRole)
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
		Data:    userData,
	}

	json.NewEncoder(w).Encode(fetchUserResponse)
}

func (api *API) lihatDataDokter(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	dokterData, err := api.usersRepo.FetchDataDokter()
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
		Data:    dokterData,
	}

	json.NewEncoder(w).Encode(fetchUserResponse)
}

func (api *API) insertPetugas(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	var petugas Petugas
	err := json.NewDecoder(req.Body).Decode(&petugas)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.adminRepo.InsertPetugas(petugas.Email, petugas.Name, petugas.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	insertUser := Petugas{
		Email: petugas.Email,
		Name:  petugas.Name,
	}

	insertPetugasResponse := InsertPetugasSuccessResponse{
		Message: "insert success",
		Data:    insertUser,
	}

	json.NewEncoder(w).Encode(insertPetugasResponse)
}
