package main

import (
	"net/http"

	"github.com/TripleTripTeam/serverV2/handler"
)

func main() {
	gbChan := make(chan handler.Coord, 10)
	h := handler.NewHandler(gbChan)

	go h.CoordProcess()

	http.HandleFunc("/moveCar", h.MoveCar)
	http.HandleFunc("/sendCoords", h.SendCoords)

	http.ListenAndServe("localhost:8000", nil)
}
