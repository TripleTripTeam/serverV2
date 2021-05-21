package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

type moveResponse struct {
	Speed float64 `json:"speed"`
	Angle float64 `json:"angle"`
}

func (h *Handler) MoveCar(w http.ResponseWriter, r *http.Request) {

	var resp moveResponse
	w.Header().Set("Content-Type", "application/json")

	resp.Speed = 10.0 - rand.Float64()
	resp.Angle = 1 - rand.Float64()*3

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}
