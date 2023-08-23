package api

import (
	"encoding/json"
	"net/http"
)

type DokterResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (api *API) lihatJadwalReservasi(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	userID := req.Context().Value("id").(int64)

	reservasi, err := api.dokterRepo.LihatJadwal(userID)
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(ReservasiErrorResponse{Error: err.Error()})
		}
	}()
	dokterResponse := DokterResponse{
		Message: "success",
		Data:    reservasi,
	}

	json.NewEncoder(w).Encode(dokterResponse)
}
