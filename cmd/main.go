package main

import (
	"net/http"

	"github.com/TripleTripTeam/serverV2/handler"
)

func main() {
	h := handler.NewHandler()

	http.HandleFunc("/moveCar", h.MoveCar)
	
	http.ListenAndServe("localhost:8000", nil)
}
