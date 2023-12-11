package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type DokterResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (api *API) lihatJadwalReservasi(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	userID := req.Context().Value("id").(int64)
	page, err := strconv.Atoi(req.URL.Query().Get("page"))

	reservasi, err := api.reservasiRepo.LihatJadwalReservasi(userID, page)
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
