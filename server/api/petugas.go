package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type VerifReservasi struct {
	ID     int64  `db:"id"`
	Status string `db:"status"`
}

type VerifikasiResponse struct {
	Message string `json:"message"`
}

type ReservasiResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
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

	fmt.Println(bodyRequest.ID)
	fmt.Println(bodyRequest.Status)

	verifikasiResponse := VerifikasiResponse{
		Message: "Berhasil merubah status reservasi",
	}

	json.NewEncoder(w).Encode(verifikasiResponse)
}
