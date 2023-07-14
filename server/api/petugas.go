package api

import (
	"encoding/json"
	"net/http"
)

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
